package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/ipuppet/aqua/server/pkg/infra/config"
	"github.com/ipuppet/aqua/server/pkg/infra/logutil"
	"github.com/ipuppet/aqua/server/pkg/infra/models"
	"github.com/ipuppet/aqua/server/pkg/infra/telemetry"
	"github.com/ipuppet/aqua/server/pkg/transport/grpcx"
	"github.com/ipuppet/aqua/server/pkg/transport/middleware"
)

// NewGRPCServer returns a gRPC server with internal auth interceptor.
// mount may be nil; callers can also register services on the returned server later.
func NewGRPCServer(rt *Runtime, mount func(server *grpc.Server)) *grpc.Server {
	if rt == nil {
		panic("bootstrap.NewGRPCServer: runtime is nil")
	}
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcx.UnaryServerInternalAuth(func() string {
				appCfg := rt.AppConfig()
				if appCfg == nil {
					return ""
				}
				return appCfg.Internal.Token
			}),
			grpcx.UnaryServerErrorMapping(),
		),
	)

	if mount != nil {
		mount(server)
	}

	return server
}

// NewHTTPServer returns a newly created HTTP server instance configured with standard middlewares.
// rt must be non-nil and contain a loaded AppConfig with Otel configuration initialized.
// mount may be nil.
func NewHTTPServer(rt *Runtime, mount func(engine *gin.Engine)) *http.Server {
	if rt == nil {
		panic("bootstrap.NewHTTPServer: runtime is nil")
	}
	appCfg := rt.AppConfig()
	if appCfg == nil {
		panic("bootstrap.NewHTTPServer: config is nil")
	}

	// 使用 gin.New() 而非 Default() 来精准控制中间件链
	// 禁用内置 Logger() 使用 OpenTelemetry trace
	engine := gin.New()

	// 使用 OpenTelemetry HTTP 链路追踪，并在 LoggerMiddleware 中挂载 traceID/spanID
	engine.Use(otelgin.Middleware(
		appCfg.Telemetry.Otel.Namespace,
		otelgin.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/ping"
		}),
	))

	// 为所有请求注入带 trace 上下文的 logger
	engine.Use(middleware.LoggerMiddleware(rt.Logger))

	// 拦截 panic (转化并汇报给 OTLP)，且处理本服务所有的错误响应格式
	engine.Use(middleware.ErrorHandler())

	// Shared health probe so mono and single-service deploys expose exactly one /ping.
	// Plain body write (not response.*): uptime probes expect literal "pang", not JSON envelopes.
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
		_, _ = ctx.Writer.Write([]byte("pang"))
	})

	if mount != nil {
		mount(engine)
	}

	// 设置较短的读写超时时间，防止慢网络客户端消耗过多服务器资源。
	// ReadTimeout 5秒可满足普通请求头的接收，WriteTimeout 10秒允许较大的响应体传输。
	return &http.Server{
		Addr:         rt.Env.HTTPAddress,
		Handler:      engine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

// startHTTPInGroup schedules the HTTP server to run inside the errgroup.
// It explicitly ignores [http.ErrServerClosed] to avoid misinterpreting graceful shutdowns as errors.
func startHTTPInGroup(g *errgroup.Group, logger *slog.Logger, srv *http.Server) {
	g.Go(func() error {
		if logger != nil {
			logger.Info("starting HTTP server", slog.String("addr", srv.Addr))
		}
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("http serve: %w", err)
		}
		return nil
	})
}

// startGRPCInGroup schedules the gRPC server to run inside the errgroup.
// addr is the listen address (e.g. env.GRPCAddr from flags).
func startGRPCInGroup(g *errgroup.Group, logger *slog.Logger, serviceName, addr string, srv *grpc.Server) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen gRPC on %s: %w", addr, err)
	}
	if logger != nil {
		logger.Info("starting gRPC server", "service", serviceName, "addr", addr)
	}
	g.Go(func() error {
		if err := srv.Serve(lis); err != nil {
			if logger != nil {
				logger.Error("gRPC server stopped", "service", serviceName, "error", err)
			}
			return fmt.Errorf("grpc serve: %w", err)
		}
		return nil
	})
	return nil
}

// Module describes how a service mounts onto shared HTTP/gRPC listeners.
// Either or both mount functions may be nil.
type Module struct {
	MountHTTP func(*gin.Engine)
	MountGRPC func(*grpc.Server)
}

// listeners holds the process-level servers started by Run / RunAll.
type listeners struct {
	HTTP *http.Server
	GRPC *grpc.Server
}

// Runtime holds process-scoped dependencies passed from Run to service assemblers.
// It is intentionally thin (env/config/logger/cleanup) — not an IoC container; Wire owns DI.
type Runtime struct {
	Env           Env
	ConfigWatcher *config.Watcher[*models.AppConfig]
	ConfigLoader  *config.Loader
	Logger        *slog.Logger
	Cleanup       *CleanupQueue
}

// AppConfig returns the latest application config snapshot from the watcher.
func (rt *Runtime) AppConfig() *models.AppConfig {
	if rt == nil || rt.ConfigWatcher == nil {
		return nil
	}
	return rt.ConfigWatcher.Value()
}

// Service defines the configuration for booting a service. Factory returns a Module;
// etcd registration is performed automatically when MountGRPC is non-nil.
type Service struct {
	Name    string
	Factory func(rt *Runtime) (*Module, error)
}

// Run boots a single service: starts HTTP and/or gRPC, registers with etcd when
// MountGRPC is present, and orchestrates graceful shutdown.
func Run(service Service) {
	runServices(service.Name, []Service{service})
}

// RunAll boots multiple services in one process: shared Runtime, one HTTP listener,
// one gRPC listener (all MountGRPC handlers), and etcd registration per gRPC service.
func RunAll(services []Service) {
	runServices("mono", services)
}

func runServices(hostName string, services []Service) {
	if len(services) == 0 {
		slog.Error("no services to run")
		os.Exit(1)
	}

	env := mustParseEnv(hostName)
	rt := mustBuildRuntime(hostName, env)
	mustInitTelemetry(hostName, rt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, _ := errgroup.WithContext(ctx)

	running := mustStartServices(hostName, services, rt, g)

	waitForCleanup := SetupGracefulShutdown(running, rt.Cleanup, rt.Logger)

	if err := g.Wait(); err != nil {
		rt.Logger.Error("server runtime error", slog.Any("error", err))
	}

	waitForCleanup()
	rt.Logger.Info("all servers stopped")
}

// mustParseEnv parses command-line flags into Env and exits the process on error.
func mustParseEnv(serviceName string) Env {
	env, err := ParseFlags()
	if err != nil {
		slog.Error("failed to parse flags",
			slog.String("service", serviceName),
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	return env
}

// mustBuildRuntime constructs the process Runtime, initializing logger and config loader.
// Any unrecoverable error during initialization will be logged and cause the process to exit.
func mustBuildRuntime(serviceName string, env Env) *Runtime {
	cleanup := NewCleanupQueue()
	logger := logutil.Default()

	loader := config.NewLoader(env.ConfigPath)
	watcher, err := config.NewWatcher[*models.AppConfig](loader, "app.yaml", logger)
	if err != nil {
		logger.Error("failed to load config",
			slog.String("service", serviceName),
			slog.String("config_path", env.ConfigPath),
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	return &Runtime{
		Env:           env,
		ConfigWatcher: watcher,
		ConfigLoader:  loader,
		Logger:        logger,
		Cleanup:       cleanup,
	}
}

// mustInitTelemetry initializes the OpenTelemetry provider for the given service.
// OTel is treated as a required dependency: initialization failure will be logged and cause the process to exit.
func mustInitTelemetry(serviceName string, rt *Runtime) {
	appCfg := rt.AppConfig()
	if appCfg == nil {
		rt.Logger.Error("missing app config for telemetry initialization",
			slog.String("service", serviceName),
		)
		os.Exit(1)
	}

	if appCfg.Telemetry.Otel.Enabled {
		otelCleanup, err := telemetry.InitProvider(
			context.Background(),
			appCfg.Telemetry.Otel.OTLP.Endpoint,
			appCfg.Telemetry.Otel.OTLP.Username,
			appCfg.Telemetry.Otel.OTLP.APIKey,
			appCfg.Telemetry.Otel.Namespace,
			serviceName,
		)
		if err != nil {
			rt.Logger.Error("failed to initialize OpenTelemetry provider",
				slog.String("service", serviceName),
				slog.Any("error", err),
			)
			os.Exit(1)
		}
		rt.Cleanup.Add(func() {
			_ = otelCleanup(context.Background())
		})
	} else {
		rt.Logger.Info("OpenTelemetry is disabled via config", slog.String("service", serviceName))
	}

	// Initialize the unified logger (File + OTel + Stdout) AFTER OTel provider is ready
	fileLogPath := ""
	if appCfg.Telemetry.FileLog.Enabled {
		fileLogPath = appCfg.Telemetry.FileLog.FilePath
	}
	logger, closer, err := logutil.NewProductionLogger(fileLogPath, serviceName)
	if err != nil {
		rt.Logger.Error("bootstrap logger initialization failed",
			slog.String("service", serviceName),
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	if closer != nil {
		rt.Cleanup.Add(func() {
			_ = closer.Close()
		})
	}

	rt.Logger = logger
	rt.ConfigWatcher.SetLogger(logger)
}

// mustStartServices invokes each service factory, merges mounts onto shared listeners,
// starts HTTP and/or gRPC, and registers gRPC services with etcd.
func mustStartServices(hostName string, services []Service, rt *Runtime, g *errgroup.Group) *listeners {
	var httpMounts []func(*gin.Engine)
	var grpcMounts []func(*grpc.Server)
	var grpcServiceNames []string

	for _, service := range services {
		mod, err := service.Factory(rt)
		if err != nil {
			rt.Logger.Error("service factory failed",
				slog.String("service", service.Name),
				slog.Any("error", err),
			)
			os.Exit(1)
		}
		if mod == nil {
			rt.Logger.Error("no module returned",
				slog.String("service", service.Name),
			)
			os.Exit(1)
		}
		if mod.MountHTTP == nil && mod.MountGRPC == nil {
			rt.Logger.Error("module has no HTTP or gRPC mount",
				slog.String("service", service.Name),
			)
			os.Exit(1)
		}

		if mod.MountHTTP != nil {
			httpMounts = append(httpMounts, mod.MountHTTP)
		}
		if mod.MountGRPC != nil {
			grpcMounts = append(grpcMounts, mod.MountGRPC)
			grpcServiceNames = append(grpcServiceNames, service.Name)
		}

		rt.Logger.Info("service assembled", slog.String("service", service.Name))
	}

	running := &listeners{}

	if len(httpMounts) > 0 {
		running.HTTP = NewHTTPServer(rt, func(engine *gin.Engine) {
			for _, mount := range httpMounts {
				mount(engine)
			}
		})
		startHTTPInGroup(g, rt.Logger, running.HTTP)
	}

	if len(grpcMounts) > 0 {
		running.GRPC = NewGRPCServer(rt, func(server *grpc.Server) {
			for _, mount := range grpcMounts {
				mount(server)
			}
		})
		if err := startGRPCInGroup(g, rt.Logger, hostName, rt.Env.GRPCAddress, running.GRPC); err != nil {
			rt.Logger.Error("failed to start gRPC server",
				slog.String("host", hostName),
				slog.String("grpc_addr", rt.Env.GRPCAddress),
				slog.Any("error", err),
			)
			os.Exit(1)
		}

		appCfg := rt.AppConfig()
		if appCfg == nil {
			rt.Logger.Error("failed to register services with etcd: config is nil",
				slog.String("host", hostName),
			)
		} else {
			for _, name := range grpcServiceNames {
				if err := RegisterService(name, rt.Env, appCfg.Etcd, rt.Cleanup, rt.Logger); err != nil {
					rt.Logger.Error("failed to register service with etcd",
						slog.String("service", name),
						slog.Any("error", err),
					)
				} else {
					rt.Logger.Info("service registered with etcd",
						slog.String("service", name),
					)
				}
			}
		}
	}

	return running
}

// SetupGracefulShutdown hooks os signals to orchestrate safe server teardown and runs remaining pending tasks in the [CleanupQueue].
// Shutdown order: HTTP first, then gRPC, then CleanupQueue.
//
// 限制最长 30s 的超时时间，避免在极端断网情况下部分连接无法关闭而导致整个服务永久挂起。
func SetupGracefulShutdown(running *listeners, cleanup *CleanupQueue, logger *slog.Logger) func() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	cleanupDone := make(chan struct{})
	var once sync.Once

	shutdown := func(trigger string) {
		once.Do(func() {
			signal.Stop(sigChan)
			if logger != nil {
				logger.Info("initiating graceful shutdown",
					slog.String("shutdown_trigger", trigger))
			}

			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer shutdownCancel()

			if running != nil && running.HTTP != nil {
				if err := running.HTTP.Shutdown(shutdownCtx); err != nil {
					if logger != nil {
						logger.Error("failed to shutdown HTTP server", slog.Any("error", err))
					}
				} else if logger != nil {
					logger.Info("HTTP server stopped")
				}
			}

			if running != nil && running.GRPC != nil {
				running.GRPC.GracefulStop()
				if logger != nil {
					logger.Info("gRPC server stopped")
				}
			}

			cleanup.Run(logger)
			close(cleanupDone)
		})
	}

	go func() {
		select {
		case sig := <-sigChan:
			shutdown(sig.String())
		case <-cleanupDone:
		}
	}()

	return func() {
		shutdown("manual")
		<-cleanupDone
	}
}

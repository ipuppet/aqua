package models

// ExtrasConfig stores any unmodeled configuration fields (mapstructure remain).
type ExtrasConfig map[string]any

// AppConfig represents the root configuration structure for the application.
type AppConfig struct {
	Database  DatabaseConfig        `mapstructure:"database" json:"database"`
	Redis     RedisConfig           `mapstructure:"redis" json:"redis"`
	Etcd      EtcdConfig            `mapstructure:"etcd" json:"etcd"`
	Email     EmailConfig           `mapstructure:"email" json:"email"`
	Turnstile TurnstileConfig       `mapstructure:"turnstile" json:"turnstile"`
	Internal  InternalServiceConfig `mapstructure:"internal_service" json:"internal_service"`
	Telemetry TelemetryConfig       `mapstructure:"telemetry" json:"telemetry"`
	Meili     MeiliSearchConfig     `mapstructure:"meilisearch" json:"meilisearch"`
	JWT       JWTConfig             `mapstructure:"jwt" json:"jwt"`
	Extras    ExtrasConfig          `mapstructure:",remain" json:"-"`
}

// MeiliSearchConfig represents the connection settings for the MeiliSearch engine.
type MeiliSearchConfig struct {
	Host   string `mapstructure:"host" json:"host"`
	APIKey string `mapstructure:"api_key" json:"api_key"`
}

// TelemetryConfig groups both file logging and Otel configurations for unified observability.
type TelemetryConfig struct {
	FileLog FileLogConfig `mapstructure:"file_log" json:"file_log"`
	Otel    OtelConfig    `mapstructure:"otel" json:"otel"`
}

// FileLogConfig represents configuration for optional local file logging.
type FileLogConfig struct {
	Enabled  bool   `mapstructure:"enabled" json:"enabled"`
	FilePath string `mapstructure:"file_path" json:"file_path"`
}

// OtelConfig represents the OpenTelemetry configuration for distributed tracing.
type OtelConfig struct {
	Enabled   bool       `mapstructure:"enabled" json:"enabled"`
	Namespace string     `mapstructure:"namespace" json:"namespace"`
	OTLP      OTLPConfig `mapstructure:"otlp" json:"otlp"`
}

// OTLPConfig holds the OpenTelemetry Protocol exporter settings.
type OTLPConfig struct {
	Endpoint string `mapstructure:"endpoint" json:"endpoint"`
	Username string `mapstructure:"username" json:"username"`
	APIKey   string `mapstructure:"api_key" json:"api_key"`
}

// DatabaseConfig represents the connection parameters for the primary relational database.
type DatabaseConfig struct {
	Driver      string `mapstructure:"driver" json:"driver"`
	Host        string `mapstructure:"host" json:"host"`
	Port        string `mapstructure:"port" json:"port"`
	Username    string `mapstructure:"username" json:"username"`
	Password    string `mapstructure:"password" json:"password"`
	Database    string `mapstructure:"database" json:"database"`
	SSLMode     string `mapstructure:"ssl_mode" json:"ssl_mode"`
	SSLRootCert string `mapstructure:"ssl_root_cert" json:"ssl_root_cert"`
	SSLCert     string `mapstructure:"ssl_cert" json:"ssl_cert"`
	SSLKey      string `mapstructure:"ssl_key" json:"ssl_key"`
}

// RedisConfig represents the connection strings and settings for the Redis cache.
type RedisConfig struct {
	Addr     string `mapstructure:"addr" json:"addr"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
}

// EtcdConfig represents the endpoint configuration for the etcd distributed key-value store.
type EtcdConfig struct {
	Endpoints []string `mapstructure:"endpoints" json:"endpoints"`
}

// EmailConfig represents the SMTP credentials and host settings for sending emails.
type EmailConfig struct {
	From     string `mapstructure:"from" json:"from"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
}

// TurnstileConfig represents the Cloudflare Turnstile verification settings.
type TurnstileConfig struct {
	SecretKey string `mapstructure:"secret_key" json:"secret_key"`
	VerifyURL string `mapstructure:"verify_url" json:"verify_url"`
}

// InternalServiceConfig represents the authentication settings for inbound gRPC.
type InternalServiceConfig struct {
	Token string `mapstructure:"token" json:"token"`
}

// JWTConfig holds RSA key paths for asymmetric JWT signing.
// The auth center (passport) holds the private key and signs tokens; other services and frontends verify with the public key.
type JWTConfig struct {
	PrivateKeyPath string `mapstructure:"private_key_path" json:"private_key_path"`
	PublicKeyPath  string `mapstructure:"public_key_path" json:"public_key_path"`
}

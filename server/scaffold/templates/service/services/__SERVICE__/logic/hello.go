package logic

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ipuppet/aqua/server/services/__SERVICE__/domain"
)

type HelloService struct {
	logger *slog.Logger
}

func NewHelloService(logger *slog.Logger) *HelloService {
	return &HelloService{logger: logger}
}

func (s *HelloService) SayHello(_ context.Context, name string) (*domain.Greeting, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "world"
	}
	msg := fmt.Sprintf("hello, %s", name)
	if s.logger != nil {
		s.logger.Debug("SayHello", slog.String("message", msg))
	}
	return &domain.Greeting{Message: msg}, nil
}

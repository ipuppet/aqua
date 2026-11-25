package grpcserver

import (
	"context"

	__SERVICE__v1 "github.com/ipuppet/aqua/server/gen/go/aqua/__SERVICE__/v1"
	"github.com/ipuppet/aqua/server/services/__SERVICE__/logic"
)

type __Service__Server struct {
	helloService *logic.HelloService
	__SERVICE__v1.Unimplemented__Service__ServiceServer
}

func New__Service__Server(helloService *logic.HelloService) *__Service__Server {
	return &__Service__Server{helloService: helloService}
}

func (s *__Service__Server) SayHello(ctx context.Context, req *__SERVICE__v1.SayHelloRequest) (*__SERVICE__v1.SayHelloResponse, error) {
	greeting, err := s.helloService.SayHello(ctx, req.GetName())
	if err != nil {
		return nil, err
	}
	return &__SERVICE__v1.SayHelloResponse{Message: greeting.Message}, nil
}

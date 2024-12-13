package grpc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"kbox-api/internal/config"
	"kbox-api/internal/shared/httperror"
	"kbox-api/pkg/proto/gengrpc"
)

type ClientConnInterface interface {
	GetRecommendations(ctx context.Context, films []*gengrpc.Film) ([]string, error)
}

type ClientConn struct {
	conn   *grpc.ClientConn
	client gengrpc.RecommendationUseCaseClient
}

func NewClient(config *config.Config) (*ClientConn, error) {
	conn, err := grpc.NewClient(config.App.GrpcServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, httperror.New(
			http.StatusInternalServerError,
			fmt.Sprintf("Не удалось подключиться к серверу: %v", err),
		)
	}

	client := gengrpc.NewRecommendationUseCaseClient(conn)
	return &ClientConn{conn: conn, client: client}, nil
}

func (c *ClientConn) Close() error {
	return c.conn.Close()
}

func (c *ClientConn) GetRecommendations(ctx context.Context, films []*gengrpc.Film) ([]string, error) {
	request := &gengrpc.RecommendationRequest{
		Films: films,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	response, err := c.client.GetRecommendations(ctx, request)
	if err != nil {
		return nil, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return response.Films, nil
}

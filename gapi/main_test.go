package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/julysNICK/simplebank/db/sqlc"
	"github.com/julysNICK/simplebank/token"
	"github.com/julysNICK/simplebank/utils"
	"github.com/julysNICK/simplebank/worker"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc/metadata"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := utils.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, duration time.Duration) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	bearerToken := fmt.Sprintf("%s %s", authorizationTypeBearer, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{
			bearerToken,
		},
	}

	return metadata.NewIncomingContext(context.Background(), md)
}
package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/julysNICK/simplebank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationTypeBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload , error){
	md,ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, fmt.Errorf("metadata is not provided")
	}

values :=	md.Get(authorizationHeader)

	if len(authorizationHeader) == 0 {
		return nil, fmt.Errorf("authorization header is not provided")
	}

	authHeader := values[0]
	
	fields := strings.Fields(authHeader)

	if len(fields) != 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])

	if authType != authorizationTypeBearer {
		return nil, fmt.Errorf("authorization type is not supported")
	}

	accessToken := fields[1]

	payload, err := server.tokenMaker.VerifyToken(accessToken)

	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	return payload, nil



}
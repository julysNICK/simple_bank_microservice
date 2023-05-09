package gapi

import (
	"fmt"

	db "github.com/julysNICK/simplebank/db/sqlc"
	"github.com/julysNICK/simplebank/token"
	"github.com/julysNICK/simplebank/utils"
)

type Server struct {
	pb.UnimplementedBankServiceServer
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new GRPC server 
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}


	return server, nil
}

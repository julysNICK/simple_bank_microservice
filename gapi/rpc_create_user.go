package gapi

import (
	"context"

	db "github.com/julysNICK/simplebank/db/sqlc"
	"github.com/julysNICK/simplebank/pb"
	"github.com/julysNICK/simplebank/utils"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hash, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot hash password: %v", err)

	}

	arg := db.CreateUserParams{
		Username:       req.GetEmail(),
		HashedPassword: hash,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.InvalidArgument, "username or email already exist")
			}
		}

		return nil, status.Errorf(codes.Internal, "cannot create user: %v", err)

	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}

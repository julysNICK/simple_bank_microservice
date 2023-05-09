package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/julysNICK/simplebank/api"
	db "github.com/julysNICK/simplebank/db/sqlc"
	"github.com/julysNICK/simplebank/gapi"
	"github.com/julysNICK/simplebank/utils"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {


	config, err := utils.LoadConfig(".") // load config from .env file

	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDrive, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)

	server, err := api.NewServer(*config, store)
	
	if err != nil {
		print("config.ServerAddress" + config.HTTPServerAddress)
		log.Fatal("cannot create server: ", err)
	}
	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

}


func runGRPCServer(config *utils.Config, store db.Store){
	server, err := gapi.NewServer(*config, store)
	if err != nil {
		print("config.ServerAddress" + config.GRPCServerAddress)
		log.Fatal("cannot create server: ", err)
	}
		grpcServer := grpc.NewServer()
		pb.RegisterBankServiceServer(grpcServer, server)
		reflection.Register(grpcServer)

		listener, err := net.Listen("tcp", config.GRPCServerAddress)

		if err != nil {
			log.Fatal("cannot start server: ", err)
		}

		log.Printf("starting gRPC server on %s", config.GRPCServerAddress)

		err = grpcServer.Serve(listener)

		if err != nil {
			log.Fatal("cannot start server: ", err)
		}
}


func runGinServer(config *utils.Config, store db.Store){
	server, err := api.NewServer(*config, store)
	
	if err != nil {
		print("config.ServerAddress" + config.HTTPServerAddress)
		log.Fatal("cannot create server: ", err)
	}
	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

}

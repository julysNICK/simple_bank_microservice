package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/julysNICK/simplebank/api"
	db "github.com/julysNICK/simplebank/db/sqlc"
	"github.com/julysNICK/simplebank/gapi"
	"github.com/julysNICK/simplebank/pb"

	"github.com/julysNICK/simplebank/utils"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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

	go runGatewayServer(config, store)
	runGRPCServer(config, store)
}

func runGRPCServer(config *utils.Config, store db.Store) {
	server, err := gapi.NewServer(*config, store)
	if err != nil {
		print("config.ServerAddress" + config.GRPCServerAddress)
		log.Fatal("cannot create server: ", err)
	}
	grpcServer := grpc.NewServer()
	// pb.RegisterBankServiceServer(grpcServer, server)
	pb.RegisterSimpleBankServer(grpcServer, server)
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

func runGatewayServer(config *utils.Config, store db.Store) {
	server, err := gapi.NewServer(*config, store)
	if err != nil {
		print("config.ServerAddress" + config.GRPCServerAddress)
		log.Fatal("cannot create server: ", err)
	}
	jsonOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOptions)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatal("cannot register handler server: ", err)
	}

	mux := http.NewServeMux()

	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.Dir("./doc/swagger"))

	mux.Handle("/swagger/", http.StripPrefix("/swagger", fs))

	listener, err := net.Listen("tcp", config.HTTPServerAddress)

	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	log.Printf("starting http server on %s", config.HTTPServerAddress)

	err = http.Serve(listener, mux)

	if err != nil {
		log.Fatal("cannot http gateway server: ", err)
	}
}

func runGinServer(config *utils.Config, store db.Store) {
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

package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/julysNICK/simplebank/api"
	db "github.com/julysNICK/simplebank/db/sqlc"
	_ "github.com/julysNICK/simplebank/doc/statik"
	"github.com/julysNICK/simplebank/gapi"
	"github.com/julysNICK/simplebank/pb"
	"github.com/rakyll/statik/fs"

	"github.com/julysNICK/simplebank/utils"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {

	config, err := utils.LoadConfig(".") // load config from .env file

	if err != nil {
		log.Fatal().Msgf("cannot load config: ", err)
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDrive, config.DBSource)

	if err != nil {
		log.Fatal().Msgf("cannot connect to db: ", err)
	}
	runDBMigration(config.MigrationUrl, config.DBSource)

	store := db.NewStore(conn)

	go runGatewayServer(config, store)
	runGRPCServer(config, store)
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)

	if err != nil {
		log.Fatal().Msgf("cannot create migration: ", err)
	}

	err = migration.Up()

	if err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msgf("cannot migrate db: ", err)
	}

	log.Print("migration completed")

}

func runGRPCServer(config *utils.Config, store db.Store) {
	server, err := gapi.NewServer(*config, store)
	if err != nil {
		print("config.ServerAddress" + config.GRPCServerAddress)
		log.Fatal().Msgf("cannot create server: ", err)
	}
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	// pb.RegisterBankServiceServer(grpcServer, server)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)

	if err != nil {
		log.Fatal().Msgf("cannot start server: ", err)
	}

	log.Printf("starting gRPC server on %s", config.GRPCServerAddress)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal().Msgf("cannot start server: ", err)
	}
}

func runGatewayServer(config *utils.Config, store db.Store) {
	server, err := gapi.NewServer(*config, store)
	if err != nil {
		print("config.ServerAddress" + config.GRPCServerAddress)
		log.Fatal().Msgf("cannot create server: ", err)
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
		log.Fatal().Msgf("cannot register handler server: ", err)
	}

	mux := http.NewServeMux()

	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()

	if err != nil {
		log.Fatal().Msgf("cannot load static files: ", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))

	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)

	if err != nil {
		log.Fatal().Msgf("cannot start server: ", err)
	}

	log.Info().Msgf("starting http server on %s", config.HTTPServerAddress)
	handle := gapi.HttpLogger(mux)
	err = http.Serve(listener, handle)

	if err != nil {
		log.Fatal().Msgf("cannot http gateway server: ", err)
	}
}

func runGinServer(config *utils.Config, store db.Store) {
	server, err := api.NewServer(*config, store)

	if err != nil {
		print("config.ServerAddress" + config.HTTPServerAddress)
		log.Fatal().Msgf("cannot create server: ", err)
	}
	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal().Msgf("cannot start server: ", err)
	}

}

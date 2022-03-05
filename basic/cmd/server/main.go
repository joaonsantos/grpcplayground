package main

import (
	"context"
	pb "grpc-log/login"
	"grpc-log/store"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const addr = ":8888"

var db store.Store

func init() {
	db = store.New()
}

type LogServer struct {
	pb.UnimplementedLogServer
}

func (s LogServer) List(ctx context.Context, void *pb.Void) (*pb.LoginList, error) {
	return db.List(), nil
}

func (s LogServer) Save(ctx context.Context, u *pb.User) (*pb.Login, error) {
	login := &pb.Login{Username: u.Name, LastLogin: timestamppb.New(time.Now())}
	db.Save(login.Username, login.GetLastLogin().AsTime())

	return login, nil
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("could not listen on %v", addr)
	}
	srv := grpc.NewServer()
	pb.RegisterLogServer(srv, newServer())

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)

	go func() {
		srv.Serve(lis)
	}()
	log.Printf("server started on %v\n", addr)

	<-s
	log.Println("server shutting down...")
	srv.GracefulStop()
}

func newServer() LogServer {
	return LogServer{}
}

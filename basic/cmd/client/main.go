package main

import (
	"context"
	"flag"
	"fmt"
	pb "grpc-log/login"
	"os"

	"google.golang.org/grpc"
)

const addr = ":8888"

func main() {
	if err := parseArgs(); err != nil {
		fmt.Fprintf(os.Stderr, "client: %v\n", err)
		os.Exit(1)
	}
}

func parseArgs() error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	flag.Parse()
	if flag.NArg() < 1 {
		return fmt.Errorf("missing subcommand: list or login")
	}
	switch command := flag.Arg(0); command {
	case "list":
		if err := list(context.Background(), pb.NewLogClient(conn)); err != nil {
			return err
		}
	case "login":
		if err := save(context.Background(), pb.NewLogClient(conn)); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized subcommand: %v", command)
	}
	return nil
}

func list(ctx context.Context, client pb.LogClient) error {
	l, err := client.List(ctx, &pb.Void{})
	if err != nil {
		return err
	}
	fmt.Printf("%+v", l)
	return nil
}

func save(ctx context.Context, client pb.LogClient) error {
	name := flag.Arg(1)
	_, err := client.Save(ctx, &pb.User{Name: name})
	if err != nil {
		return err
	}
	return nil
}

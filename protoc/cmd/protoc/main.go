package main

import (
	"flag"
	"fmt"
	"os"
	"protoc-log/login"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const logfile = "log.pb"

func main() {
	if err := parseArgs(); err != nil {
		fmt.Fprintf(os.Stderr, "greet: %v\n", err)
		os.Exit(1)
	}
}

func parseArgs() error {
	flag.Parse()
	if flag.NArg() < 1 {
		return fmt.Errorf("missing subcommand: login or last")
	}
	switch command := flag.Arg(0); command {
	case "last":
		if err := last(); err != nil {
			return err
		}
	case "login":
		if err := save(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized subcommand: %v", command)
	}
	return nil
}

func last() error {
	f, err := os.ReadFile(logfile)
	if err != nil {
		return err
	}
	login := &login.Login{}
	if err := proto.Unmarshal(f, login); err != nil {
		return err
	}
	lastLogin := login.LastLogin.AsTime().Format(time.ANSIC)
	fmt.Printf("user %q has last logged in on %v\n", login.Username, lastLogin)
	return nil
}

func save() error {
	l := login.Login{
		Username:  "user",
		LastLogin: timestamppb.New(time.Now()),
	}
	bs, err := proto.Marshal(&l)
	if err != nil {
		return err
	}
	if err := os.WriteFile(logfile, bs, 0644); err != nil {
		return err
	}
	return nil
}

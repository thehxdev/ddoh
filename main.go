package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/thehxdev/ddoh/config"
	"github.com/thehxdev/ddoh/server"
)

const (
	VERSION = "1.0.0"
)

var (
	confPath string
)

func main() {
	err := os.Setenv("GOGC", "20")
	if err != nil {
		log.Fatal(err)
	}

	configureCmdFlags()

	config.InitConfig(confPath)
	server := server.Init()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	serverCtx, serverCtxStop := context.WithCancel(context.Background())

	go func() {
		<-sigChan
		log.Println("Shutting down the server...")
		server.Shutdown()
		serverCtxStop()
	}()

	server.Start()
	<-serverCtx.Done()
}

func configureCmdFlags() {
	flag.StringVar(&confPath, "c", "config.json", "path to config.json file")
	flag.Parse()
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/thehxdev/ddoh/config"
	"github.com/thehxdev/ddoh/server"
)

const (
	VERSION = "1.0.4"
)

var (
	confPath    string
	showVersion bool
)

func main() {
	err := os.Setenv("GOGC", "20")
	if err != nil {
		log.Fatal(err)
	}

	configureCmdFlags()

	if showVersion {
		fmt.Println("ddoh v" + VERSION + "\nhttps://github.com/thehxdev/ddoh")
		os.Exit(0)
	}

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
	flag.BoolVar(&showVersion, "v", false, "show version info")
	flag.Parse()
}

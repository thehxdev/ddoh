package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/thehxdev/ddoh/config"
	"github.com/thehxdev/ddoh/server"
)

const (
	VERSION = "1.0.7"
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
		fmt.Printf("ddoh v%s\nhttps://github.com/thehxdev/ddoh", VERSION)
		os.Exit(0)
	}

	config.InitConfig(confPath)
	server := server.Init()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	serverCtx, serverCtxStop := context.WithCancel(context.Background())
	server.Ctx = serverCtx

	go func() {
		<-sigChan
		log.Println("Shutting down the server...")
		server.Shutdown()
		serverCtxStop()
	}()

	go func() {
		stat := &runtime.MemStats{}
		for {
			runtime.ReadMemStats(stat)
			log.Printf("Heap Allocations: %d KB", stat.HeapAlloc/1024)
			time.Sleep(time.Second * 5)
		}
	}()

	go func() {
		server.Start()
	}()

	<-serverCtx.Done()
}

func configureCmdFlags() {
	flag.StringVar(&confPath, "c", "config.json", "path to config.json file")
	flag.BoolVar(&showVersion, "v", false, "show version info")
	flag.Parse()
}

package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	LocalResolver string `json:"local_resolver"`
	DoHServer     string `json:"doh_server"`
	DoHIP         string `json:"doh_ip"`
	UDPBuffSize   int    `json:"udp_buffer_size"`
}

var (
	Global *Config
)

func InitConfig(path string) {
	// Default config
	Global = &Config{
		LocalResolver: "9.9.9.9",
		DoHServer:     "https://max.rethinkdns.com/dns-query",
		DoHIP:         "137.66.7.89",
		UDPBuffSize:   512,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Println("could not read config file. using default options...")
		return
	}

	if err := json.Unmarshal(data, Global); err != nil {
		log.Fatal(err)
	}
}

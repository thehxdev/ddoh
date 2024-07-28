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
	c := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(data, c); err != nil {
		log.Fatal(err)
	}

	Global = c
}

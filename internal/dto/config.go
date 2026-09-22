package dto

import "os"

type Config struct {
	FAAKey    string
	FAASecret string
}

func NewConfig() Config {
	return Config{
		FAAKey:    os.Getenv("FAA_KEY"),
		FAASecret: os.Getenv("FAA_SECRET"),
	}
}

func (c Config) NotamEnabled() bool {
	return c.FAAKey != "" && c.FAASecret != ""
}

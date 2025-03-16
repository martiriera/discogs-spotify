package main

import (
	"log"
	"os"
)

func AssertEnvVar(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("environment variable %s is required", name)
	}
	return value
}

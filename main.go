package main

import (
	"golang-course-registration/config"
	"golang-course-registration/infrastructure/server"
)

func main() {
	cfg := config.Load()

	srv := server.New(cfg)
	srv.Init()
	srv.Start()
}

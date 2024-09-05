package app

import (
	"fmt"

	"github.com/tritonol/gophmart.git/internal/config"
)

func Run() {
	// TODO: init logger

	// TODO: load config
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: init db and run migrations
}

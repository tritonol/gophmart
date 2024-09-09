package app

import (
	"context"
	"fmt"

	"github.com/tritonol/gophmart.git/internal/config"
	repoUser "github.com/tritonol/gophmart.git/internal/repository/User"
	"github.com/tritonol/gophmart.git/internal/storage/pg"
)

func Run() {
	ctx := context.Background()
	// TODO: init logger

	// TODO: add logging
	cfg, err := config.MustLoad()
	if err != nil {
		panic("can't init config")
	}

	// TODO: add logging
	db, err := pg.NewConnection(ctx, cfg.DbUri)
	if err != nil {
		fmt.Println(err)
		panic("can't init database")
	}

	// TODO: init repos
	repoAuth := repoUser.New(ctx, db)
	_ = repoAuth

	// TODO: init usecases
}

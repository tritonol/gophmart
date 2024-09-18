package app

import (
	"context"
	"fmt"

	"github.com/tritonol/gophmart.git/internal/config"
	httpserver "github.com/tritonol/gophmart.git/internal/http"
	apiAccrual "github.com/tritonol/gophmart.git/internal/repository/Accrual/api"
	repoBalance "github.com/tritonol/gophmart.git/internal/repository/Balance"
	repoOrder "github.com/tritonol/gophmart.git/internal/repository/Order"
	repoUser "github.com/tritonol/gophmart.git/internal/repository/User"
	"github.com/tritonol/gophmart.git/internal/storage/pg"
	"github.com/tritonol/gophmart.git/internal/usecase/accrual"
	"github.com/tritonol/gophmart.git/internal/usecase/auth"
	"github.com/tritonol/gophmart.git/internal/usecase/balance"
	"github.com/tritonol/gophmart.git/internal/usecase/orders"
)

func Run() {
	ctx := context.Background()

	cfg, err := config.MustLoad()
	if err != nil {
		panic("can't init config")
	}

	db, err := pg.NewConnection(ctx, cfg.DbUri)
	if err != nil {
		fmt.Println(err)
		panic("can't init database")
	}

	repoAuth := repoUser.New(ctx, db)
	repoOrder := repoOrder.New(ctx, db)
	repoBalance := repoBalance.New(ctx, db)

	repoAccrual := apiAccrual.New(cfg.AccrualAddress)

	authUc := auth.New(repoAuth)
	ordersUc := orders.New(repoOrder)
	accruarUc := accrual.New(repoOrder, repoAccrual, repoBalance)
	balanceUc := balance.New(repoBalance)
	// init server

	accruarUc.StartProcessingAccruals(ctx)

	httpserver := httpserver.New(cfg, authUc, ordersUc, balanceUc)

	httpserver.Run(ctx)
}

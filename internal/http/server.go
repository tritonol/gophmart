package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tritonol/gophmart.git/internal/config"
	"github.com/tritonol/gophmart.git/internal/models/balance"
	"github.com/tritonol/gophmart.git/internal/models/order"
	"github.com/tritonol/gophmart.git/internal/models/user"
)

type Server struct {
	r       *chi.Mux
	auth    ucAuth
	order   ucOrder
	balance ucBalance
	cfg     *config.Config
}

type ucAuth interface {
	Register(ctx context.Context, credetials user.UserCredentials) (string, error)
	Login(ctx context.Context, credentials user.UserCredentials) (string, error)
	ValidateToken(token string) (user.UserID, error)
}

type ucOrder interface {
	CreateOrder(ctx context.Context, number int64, userId user.UserID) error
	GetUserOrders(ctx context.Context, userId user.UserID) ([]*order.Order, error)
}

type ucBalance interface {
	GetBalance(ctx context.Context, userId user.UserID) (*balance.Balance, error)
	WriteOff(ctx context.Context, userId user.UserID, orderNum int64, value float64) error
	WithdrawalsHistory(ctx context.Context, userId user.UserID) ([]*balance.Transaction, error)
}

func New(cfg *config.Config, auth ucAuth, order ucOrder, balance ucBalance) *Server {
	s := &Server{
		auth:    auth,
		order:   order,
		balance: balance,
		cfg:     cfg,
	}

	r := chi.NewRouter()
	s.r = r

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	s.UserRoutes()

	return s
}

func (s *Server) UserRoutes() *chi.Mux {
	userRouter := chi.NewRouter()

	userRouter.Group(func(r chi.Router) {
		r.Post("/register", s.RegisterUser)
		r.Post("/login", s.LoginUser)
	})

	userRouter.Group(func(r chi.Router) {
		r.Use(s.Auth)

		r.Post("/orders", s.GetOrders)
		r.Get("/orders", s.CreateOrder)

		r.Get("/balance", s.GetBalance)
		r.Post("/balance/withdraw", s.WriteOff)

		r.Get("/withdrawals", s.WithdrawalsHistory)
	})

	return userRouter
}

func (s *Server) Run(ctx context.Context) {
	server := &http.Server{
		Addr:    s.cfg.RunAddress,
		Handler: s.r,
	}

	serverCtx, serverCancel := context.WithCancel(ctx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 10*time.Second)
		defer shutdownCancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				panic("shutdown timed out")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			panic(err)
		}

		serverCancel()
	}()

	err := server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}

	<-serverCtx.Done()
}

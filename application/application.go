package main

import (
	api "authenticator/internal"
	"authenticator/internal/databases/postgresql"
	"authenticator/internal/spec"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phenpessoa/gutils/netutils/httputils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println("goodbye :)")
}

func run(ctx context.Context) error {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	defer func() { _ = logger.Sync() }()

	pool, err := pgxpool.New(ctx, fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("POSTGRESQL_USER"),
		os.Getenv("POSTGRESQL_PASS"),
		os.Getenv("POSTGRESQL_HOST"),
		os.Getenv("POSTGRESQL_PORT"),
		os.Getenv("POSTGRESQL_NAME"),
	))
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return err
	}

	store := postgresql.New(pool)
	_, err = store.GetApplicationByName(ctx, "authenticator")
	if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            store.InsertApplication(ctx, "authenticator")
			fmt.Println(" \033[0;32m✔\033[0m authenticator application inserted")
        } else {
			return err
		}
    }
	_, err = store.GetUser(ctx, os.Getenv("ROOT_MAIL"))
	if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
			hash, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ROOT_PASS")), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
            store.InsertUser(ctx, postgresql.InsertUserParams{
				Email: os.Getenv("ROOT_MAIL"),
				Name: "root",
				Password: string(hash),
			})
			fmt.Println(" \033[0;32m✔\033[0m root user inserted")
        } else {
			return err
		}
	}


	r := chi.NewMux()
	r.Use(middleware.RequestID, middleware.Recoverer, httputils.ChiLogger(logger))

	si := api.NewAPI(
		pool,
		logger,
	)

	r.Mount("/", spec.Handler(&si))
	r.Handle("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("swagger"))))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	defer func() {
		const timeout = 30 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("failed to shutdown server", zap.Error(err))
		}
	}()

	errChan := make(chan error, 1)

	go func() {
		fmt.Println(" \033[0;32m✔\033[0m Server started at", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return nil
}
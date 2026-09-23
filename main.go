package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"banking/config"
	"banking/database"
	"banking/handlers"
	"banking/middleware"
	"banking/services/account"
	"banking/services/ledger"
	"banking/services/loan"
)

func main() {
	cfg := config.Load()
	if cfg.APIKey == "" {
		log.Fatal("API_KEY must be set (empty key would accept any request)")
	}

	db, err := database.NewSQLiteDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer db.Close()

	accountRepo := database.NewAccountRepository(db)
	ledgerRepo := database.NewLedgerRepository(db)
	loanRepo := database.NewLoanRepository(db)

	ledgerSvc := ledger.NewLedgerService(ledgerRepo)
	accountSvc := account.NewAccountService(accountRepo, ledgerSvc, db)
	loanSvc := loan.NewLoanService(loanRepo, accountSvc)

	h := handlers.NewHandler(accountSvc, loanSvc, ledgerSvc)
	router := handlers.NewRouter(h, middleware.APIKeyAuth(cfg.APIKey), middleware.RateLimit(), cfg.FrontendOrigin)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down, finishing in-flight requests...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("forced shutdown: %v", err)
	}
	log.Println("server stopped")
}

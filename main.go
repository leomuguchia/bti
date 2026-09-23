package main

import (
	"log"

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
	router := handlers.NewRouter(h, middleware.APIKeyAuth(cfg.APIKey), middleware.RateLimit())

	log.Printf("listening on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

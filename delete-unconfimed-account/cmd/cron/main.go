package main

import (
	"context"
	"delete-unconfirmed-account/internal/configuration"
	"delete-unconfirmed-account/internal/database"
	"delete-unconfirmed-account/internal/job"
	"delete-unconfirmed-account/pkg/account"
	"log"
	"time"
)

const (
	START_TIME     = 5
	DELAY_TIME     = 30
	NUM_GOROUTINES = 3
)

func main() {
	config := configuration.LoadConfig()
	configDatabase := config.BuildDatabaseConfig()
	db := database.NewDatabase(configDatabase)
	accountRepository := account.NewAccountRepository(db)
	accountService := account.NewAccountService(accountRepository)

	for i := 0; i < NUM_GOROUTINES; i++ {
		go func(id int) {
			job.Cron(context.Background(), START_TIME*time.Second, DELAY_TIME*time.Second, func() {
				accountService.DeleteUnconfimed()
				log.Printf("Rodando go routine %d\n", id)
			})
		}(i)
	}

	for {
	}
}

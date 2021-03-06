package main

import (
	"context"
	"delete-unconfirmed-account/internal/cache"
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
	configRedis := config.BuildRedisConfig()
	db := database.NewDatabase(configDatabase)
	accountRepository := account.NewAccountRepository(db)
	cache := cache.NewRedisClient(configRedis)
	accountService := account.NewAccountService(accountRepository, cache)

	for i := 0; i < NUM_GOROUTINES; i++ {
		go func(id int) {
			job.Cron(context.Background(), START_TIME*time.Second, DELAY_TIME*time.Second, func() {
				accountService.DeleteUnconfimed()
				log.Printf("Rodando go routine %d\n", id)
			})
		}(i)
	}
	// just to be running until some files are proessed
	select {}
}

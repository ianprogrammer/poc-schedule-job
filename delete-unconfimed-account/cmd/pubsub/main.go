package main

import (
	"context"
	"delete-unconfirmed-account/internal/cache"
	"delete-unconfirmed-account/internal/configuration"
	"delete-unconfirmed-account/internal/database"
	"delete-unconfirmed-account/pkg/account"
)

func main() {
	config := configuration.LoadConfig()
	configDatabase := config.BuildDatabaseConfig()
	configRedis := config.BuildRedisConfig()
	db := database.NewDatabase(configDatabase)
	accountRepository := account.NewAccountRepository(db)
	cache := cache.NewRedisClient(configRedis)
	accountService := account.NewAccountService(accountRepository, cache)

	accountService.WatchExpirationEvent(context.Background())
}

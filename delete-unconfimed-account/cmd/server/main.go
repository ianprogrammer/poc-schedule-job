package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"delete-unconfirmed-account/internal/configuration"
	"delete-unconfirmed-account/internal/database"
	"delete-unconfirmed-account/pkg/account"

	"github.com/labstack/echo"
)

type App struct{}

func (app *App) Run() error {

	config := configuration.LoadConfig()

	configServer := config.BuildServerConfig()
	configDatabase := config.BuildDatabaseConfig()
	db := database.NewDatabase(configDatabase)
	accountRepository := account.NewAccountRepository(db)
	accountService := account.NewAccountService(accountRepository)

	e := echo.New()

	account.RegisterProductHandlers(e, accountService)

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", configServer.Port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("desligando server", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configServer.GracefullShutdownTimeout)*time.Second)

	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	return nil
}

func main() {
	fmt.Println("Iniciando POC accounts API")
	app := App{}

	if err := app.Run(); err != nil {
		fmt.Println(err)
	}
}

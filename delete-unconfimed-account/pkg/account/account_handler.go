package account

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type AccountRequest struct {
	Name string `json:"name"`
}

type Handler struct {
	E              *echo.Echo
	AccountService IAccountService
}

func RegisterProductHandlers(e *echo.Echo, accountService IAccountService) {
	ah := &Handler{
		E:              e,
		AccountService: accountService,
	}
	ah.E.POST("/accounts", ah.saveAccount)
	ah.E.POST("/accounts/:id/confirm", ah.accountConfirm)
	ah.E.DELETE("/accounts/:id", ah.deleteProduct)
}

func (ah *Handler) saveAccount(c echo.Context) error {
	a := new(AccountRequest)
	if err := c.Bind(a); err != nil {
		return err
	}

	result, err := ah.AccountService.Save(Account{
		Name:                 a.Name,
		RequiredConfirmation: true,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (ad *Handler) accountConfirm(c echo.Context) error {
	paramId := c.Param("id")
	id, _ := strconv.Atoi(paramId)
	p := new(AccountRequest)

	if err := c.Bind(p); err != nil {
		return err
	}

	product, err := ad.AccountService.ConfirmAccount(id, Account{
		Name:                 p.Name,
		ID:                   id,
		RequiredConfirmation: false,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, product)
}

func (ah *Handler) deleteProduct(c echo.Context) error {
	paramId := c.Param("id")
	id, _ := strconv.Atoi(paramId)
	err := ah.AccountService.DeleteById(id)

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

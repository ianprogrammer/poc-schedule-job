package account

import (
	"fmt"
	"log"
)

type AccountService struct {
	accountRepository IAccountRepository
}

type IAccountService interface {
	DeleteUnconfimed()
}

func NewAccountService(accountRepository IAccountRepository) *AccountService {
	return &AccountService{
		accountRepository: accountRepository,
	}
}

func (as *AccountService) DeleteUnconfimed() {
	err := as.accountRepository.DeleteUnconfimed()
	if err != nil {
		log.Fatal(fmt.Errorf("não foi possível deletar os usuários com emails não confirmados, %s", err.Error()))
	}
	log.Println("contas não confirmadas deletadas")
}

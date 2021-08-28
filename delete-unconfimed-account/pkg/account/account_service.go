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
	Save(a Account) (Account, error)
	DeleteById(id int) error
	ConfirmAccount(id int, a Account) (Account, error)
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

func (as *AccountService) Save(a Account) (Account, error) {
	result, err := as.accountRepository.Insert(a)
	if err != nil {
		log.Fatal("não foi possível salvar a conta")
		return Account{}, err
	}
	return result, nil
}

func (as *AccountService) DeleteById(id int) error {
	return as.accountRepository.DeleteById(id)
}

func (as *AccountService) ConfirmAccount(id int, a Account) (Account, error) {
	a.RequiredConfirmation = false
	result, err := as.accountRepository.Update(id, a)
	if err != nil {
		log.Fatal("não foi possível confirmar a conta")
		return Account{}, err
	}
	return result, nil
}

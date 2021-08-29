package account

import (
	"context"
	"delete-unconfirmed-account/internal/cache"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type AccountService struct {
	accountRepository IAccountRepository
	cacheClient       *cache.CacheClient
}

type IAccountService interface {
	DeleteUnconfimed()
	Save(a Account) (Account, error)
	DeleteById(id int) error
	ConfirmAccount(id int, a Account) (Account, error)
}

const (
	ConfirmAccountPrefixKey       = "email:confirmation"
	ConfirmationAccountExpiration = 60 * time.Second
)

func NewAccountService(accountRepository IAccountRepository, cache *cache.CacheClient) *AccountService {
	return &AccountService{
		accountRepository: accountRepository,
		cacheClient:       cache,
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

	b, _ := json.Marshal(result)
	err = as.cacheClient.Add(fmt.Sprintf("%s-%d", ConfirmAccountPrefixKey, result.ID), string(b), ConfirmationAccountExpiration)
	if err != nil {
		log.Fatal("não foi possível criar registro no redis")
		return Account{}, err
	}

	return result, nil
}

func (as *AccountService) DeleteById(id int) error {
	return as.accountRepository.DeleteById(id)
}

func (as *AccountService) ConfirmAccount(id int, a Account) (Account, error) {
	result, err := as.accountRepository.ConfirmAccount(id)
	if err != nil {
		log.Fatal("não foi possível confirmar a conta")
		return Account{}, err
	}
	as.cacheClient.RedisClient.Del(fmt.Sprintf("%s-%d", ConfirmAccountPrefixKey, result.ID))
	return result, nil
}

func (as *AccountService) WatchExpirationEvent(ctx context.Context) {

	ps := as.cacheClient.RedisClient.PSubscribe("__key*__:expired")
	ch := ps.Channel()
	for {
		select {
		case i := <-ch:
			fmt.Printf("deletando key: %s", i.Payload)
			as.expiredConfirmation(i.Payload)
		case <-ctx.Done():
			fmt.Println("ok")
		}
	}
}

func (as *AccountService) expiredConfirmation(key string) {
	s := strings.Split(key, "-")
	keyId := s[len(s)-1]
	id, _ := strconv.Atoi(keyId)
	if err := as.DeleteById(id); err != nil {
		fmt.Println(err)
	}
}

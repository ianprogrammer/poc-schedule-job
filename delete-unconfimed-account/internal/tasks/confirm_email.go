package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"delete-unconfirmed-account/pkg/account"

	"github.com/hibiken/asynq"
)

const (
	TypeEmailConfirm = "email:confirm"
)

type EmailConfirmPayload struct {
	UserID int
}

func NewEmailConfirmTask(userId int) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailConfirmPayload{UserID: userId})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailConfirm, payload), nil
}

func HandleEmailConfirmTask(_ context.Context, t *asynq.Task, accountService account.IAccountService) error {
	var p EmailConfirmPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed %v: %w", err, asynq.SkipRetry)
	}
	if err := accountService.DeleteById(p.UserID); err != nil {
		return fmt.Errorf("could be possible to delete this account")
	}
	return nil
}

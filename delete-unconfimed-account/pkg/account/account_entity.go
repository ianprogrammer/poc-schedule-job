package account

import "time"

type Account struct {
	ID                   int `gorm:"primaryKey"`
	Name                 string
	CreatedAt            *time.Time `gorm:"created_at"`
	RequiredConfirmation bool       `gorm:"required_confirmation"`
}

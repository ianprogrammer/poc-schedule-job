package account

import (
	"time"

	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

type IAccountRepository interface {
	DeleteUnconfimed() error
	Insert(createdAt time.Time, requiredConfirmation bool) (Account, error)
	DeleteById(id int) error
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) Insert(createdAt time.Time, requiredConfirmation bool) (Account, error) {
	result := Account{
		CreatedAt:            &createdAt,
		RequiredConfirmation: requiredConfirmation,
	}
	if err := r.db.Save(&result); err.Error != nil {
		return Account{}, err.Error
	}
	return result, nil
}

func (r *AccountRepository) DeleteById(id int) error {
	if result := r.db.Delete(&Account{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *AccountRepository) DeleteUnconfimed() error {
	// begin a transaction
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	auroraDBIsolationLevelQuery := "SET TRANSACTION ISOLATION LEVEL REPEATABLE READ"
	if err := tx.Exec(auroraDBIsolationLevelQuery).Error; err != nil {
		tx.Rollback()
		return err
	}

	moveQuery := `INSERT INTO users_deleted 
				  SELECT * FROM USERS WHERE created_at < now() - INTERVAL '1 day' AND required_confirmation = true`
	if err := tx.Exec(moveQuery).Error; err != nil {
		tx.Rollback()
		return err
	}

	deleteQuery := `DELETE FROM public.users 
						  WHERE id IN (SELECT id FROM users WHERE created_at < now() - interval '1 day' AND 
									   required_confirmation = true)`
	if err := tx.Exec(deleteQuery).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

}

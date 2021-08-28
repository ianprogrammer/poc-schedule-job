package account

import (
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

type IAccountRepository interface {
	DeleteUnconfimedAccounts() error
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) DeleteUnconfimedAccounts() error {
	// begin a transaction
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// auroraDBIsolationLevelQuery := "SET TRANSACTION ISOLATION LEVEL REPEATABLE READ"
	// if err := tx.Exec(auroraDBIsolationLevelQuery).Error; err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

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

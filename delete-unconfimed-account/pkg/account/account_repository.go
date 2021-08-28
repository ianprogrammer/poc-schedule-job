package account

import (
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

type IAccountRepository interface {
	DeleteUnconfimed() error
	Insert(account Account) (Account, error)
	DeleteById(id int) error
	SelectById(id int) (Account, error)
	Update(id int, a Account) (Account, error)
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) Insert(account Account) (Account, error) {
	if err := r.db.Save(&account); err.Error != nil {
		return Account{}, err.Error
	}
	return account, nil
}

func (r *AccountRepository) DeleteById(id int) error {
	if result := r.db.Delete(&Account{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *AccountRepository) SelectById(id int) (Account, error) {
	var account Account

	if result := r.db.First(&account, id); result.Error != nil {
		return account, result.Error
	}

	return account, nil

}

func (r *AccountRepository) Update(id int, a Account) (Account, error) {
	account, err := r.SelectById(id)

	if err != nil {
		return Account{}, err
	}

	if result := r.db.Model(&account).Updates(a); result.Error != nil {
		return Account{}, result.Error
	}
	return account, nil
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

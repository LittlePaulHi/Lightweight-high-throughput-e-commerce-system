package mariadb

import (
	"errors"

	"gorm.io/gorm"
)

// add index at needed field(s)
type Account struct {
	ID     int    `gorm:"primaryKey;autoIncrement;uniqueIndex"`
	Name   string `gorm:"size:100;not null;unique"`
	Amount int    `gorm:"not null"`
	Cart   Cart   `gorm:"ForeignKey:AccountID"`
	Order  Order  `gorm:"ForeignKey:AccountID"`
}

func (account *Account) Initialize(name string, amount int) {
	account.ID = 0
	account.Name = name
	account.Amount = amount
}

func (account *Account) SaveAccount() (*Account, error) {
	err := db.Create(&account).Error
	if err != nil {
		return &Account{}, err
	}
	return account, nil
}

func (account *Account) UpdateAccount(id int) (*Account, error) {
	tx := db.Model(&Account{}).Where("ID = ?", id).Take(&Account{}).UpdateColumns(
		map[string]interface{}{
			"Amount": account.Amount,
		},
	)
	if tx.Error != nil {
		return &Account{}, tx.Error
	}

	// check the updated account
	err := tx.Model(&Account{}).Where("ID = ?", id).Take(&account).Error
	if err != nil {
		return &Account{}, err
	}

	return account, nil
}

func (account *Account) FindAccountByID(id int) (*Account, error) {
	err := db.Model(Account{}).Where("ID = ?", id).Take(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Account{}, errors.New("Account Not Found")
	}
	if err != nil {
		return &Account{}, err
	}

	return account, nil
}

func (account *Account) DeleteAccount(id int) (int64, error) {
	tx := db.Model(&Account{}).Where("ID = ?", id).Take(&Account{}).Delete(&Account{})
	if tx.Error != nil {
		return 0, tx.Error
	}

	return tx.RowsAffected, nil
}

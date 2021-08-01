package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type YChaseModel struct {
	WalletAddress string         `gorm:"primary_key" json:"walletAddress"`
	Assets        pq.StringArray `gorm:"not null;type:varchar(64)[]" json:"assets"`
}

func (u *YChaseModel) Create(db *gorm.DB) (*YChaseModel, error) {

	err := db.Debug().Create(&u).Error
	if err != nil {
		return &YChaseModel{}, err
	}
	return u, nil
}

func (u *YChaseModel) Read(db *gorm.DB, WalletAddress string) (*YChaseModel, error) {
	err := db.Debug().Model(YChaseModel{}).Where("wallet_address = ?", WalletAddress).Take(&u).Error
	if err != nil {
		return &YChaseModel{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &YChaseModel{}, errors.New("YChaseModel Not Found")
	}
	return u, err
}

func (u *YChaseModel) Update(db *gorm.DB, WalletAddress string) (*YChaseModel, error) {

	db = db.Debug().Model(&YChaseModel{}).Where("wallet_address = ?", WalletAddress).Take(&YChaseModel{}).UpdateColumns(
		map[string]interface{}{
			"assets": u.Assets,
		},
	)
	if db.Error != nil {
		return &YChaseModel{}, db.Error
	}

	return u, nil
}

func (u *YChaseModel) Delete(db *gorm.DB, WalletAddress string) (int64, error) {

	db = db.Debug().Model(&YChaseModel{}).Where("wallet_address = ?", WalletAddress).Take(&YChaseModel{}).Delete(&YChaseModel{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

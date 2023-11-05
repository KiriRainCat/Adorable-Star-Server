package dao

import (
	"adorable-star/model"

	"gorm.io/gorm"
)

var Jupiter = &JupiterDAO{DB}

type JupiterDAO struct {
	db *gorm.DB
}

func (*JupiterDAO) GetDataByUID(uid int) (*model.JupiterData, error) {
	var data model.JupiterData

	err := DB.Where("uid = ?", uid).First(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

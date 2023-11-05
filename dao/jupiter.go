package dao

import (
	"adorable-star/model"

	"gorm.io/gorm"
)

type JupiterDAO struct {
	db *gorm.DB
}

func NewJupiterDAO(db *gorm.DB) *JupiterDAO {
	return &JupiterDAO{db}
}

func (dao *JupiterDAO) GetDataByUID(uid int) (*model.JupiterData, error) {
	var data model.JupiterData

	err := dao.db.Where("uid = ?", uid).First(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

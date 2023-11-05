package dao

import (
	"adorable-star/model"
)

var Jupiter = &JupiterDAO{}

type JupiterDAO struct{}

func (*JupiterDAO) GetDataByUID(uid int) (*model.JupiterData, error) {
	var data model.JupiterData

	err := DB.Where("uid = ?", uid).First(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

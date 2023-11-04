package service

import (
	"gorm.io/gorm"
)

func NewDataService(db *gorm.DB) *DataService {
	return &DataService{db}
}

type DataService struct {
	db *gorm.DB
}

// TODO: 添加数据与数据库的对比，Message的生成等

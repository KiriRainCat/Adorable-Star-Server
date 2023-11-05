package dao

import "adorable-star/model"

type MessageDAO struct{}

func (*MessageDAO) Insert(msg model.Message) error {
	return DB.Create(msg).Error
}

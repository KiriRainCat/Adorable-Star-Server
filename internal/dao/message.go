package dao

import "adorable-star/internal/model"

var Message = &MessageDAO{}

type MessageDAO struct{}

func (*MessageDAO) Insert(msg *model.Message) error {
	return DB.Create(msg).Error
}

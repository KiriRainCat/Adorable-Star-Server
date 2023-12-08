package dao

import "adorable-star/internal/model"

var Message = &MessageDAO{}

type MessageDAO struct{}

func (*MessageDAO) Insert(msg *model.Message) error {
	return DB.Create(msg).Error
}

func (*MessageDAO) GetByID(id int) (message *model.Message, err error) {
	err = DB.First(&message, id).Error
	return
}

func (*MessageDAO) GetListByUID(uid int) (messages []*model.Message, err error) {
	err = DB.Order("created_at DESC").Find(&messages, "uid = ?", uid).Error
	return
}

func (*MessageDAO) DeleteAll(uid int) error {
	return DB.Delete(&model.Message{}, "uid = ?", uid).Error
}

func (*MessageDAO) Delete(id int) error {
	return DB.Delete(&model.Message{}, id).Error
}

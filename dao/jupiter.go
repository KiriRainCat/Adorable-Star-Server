package dao

import (
	"adorable-star/model"
	"time"
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

func (*JupiterDAO) GetCoursesByUID(uid int) (courses []*model.Course, err error) {
	err = DB.Find(&courses, "uid = ?", uid).Error
	return
}

func (*JupiterDAO) GetAssignmentsByCourseAndUID(uid int, courseTitle string) (assignments []*model.Assignment, err error) {
	err = DB.Find(&assignments, "uid = ? AND `from` = ?", uid, courseTitle).Error
	return
}

func (*JupiterDAO) PutCourse(course *model.Course) error {
	return DB.Save(course).Error
}

func (*JupiterDAO) PutAssignment(assignment *model.Assignment) error {
	return DB.Save(assignment).Error
}

func (*JupiterDAO) UpdateFetchTimeAndGPA(uid int, gpa string) error {
	return DB.Model(&model.JupiterData{}).Where("uid = ?", uid).Updates(map[string]any{"gpa": gpa, "fetched_time": time.Now()}).Error
}

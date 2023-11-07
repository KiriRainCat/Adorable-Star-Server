package dao

import (
	"adorable-star/internal/model"
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

func (*JupiterDAO) InsertCourse(course *model.Course) error {
	return DB.Create(course).Error
}

func (*JupiterDAO) InsertAssignment(assignment *model.Assignment) error {
	return DB.Create(assignment).Error
}

func (*JupiterDAO) UpdateCourse(old *model.Course, course *model.Course) error {
	// Select * to select all columns, because status update can use 0 {Default will not update for gorm}
	return DB.Model(old).Select("*").Updates(course).Error
}

func (*JupiterDAO) UpdateAssignment(old *model.Assignment, assignment *model.Assignment) error {
	// Select * to select all columns, because status update can use 0 {Default will not update for gorm}
	return DB.Model(old).Select("*").Updates(assignment).Error
}

func (*JupiterDAO) UpdateFetchTimeAndGPA(uid int, gpa string) error {
	if gpa == "" {
		return DB.Model(&model.JupiterData{}).Where("uid = ?", uid).Update("fetched_at", time.Now()).Error
	}
	return DB.Model(&model.JupiterData{}).Where("uid = ?", uid).Updates(map[string]any{"gpa": gpa, "fetched_at": time.Now()}).Error

}

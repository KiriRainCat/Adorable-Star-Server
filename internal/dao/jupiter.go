package dao

import (
	"adorable-star/internal/model"
	"time"
)

var Jupiter = &JupiterDAO{}

type JupiterDAO struct{}

func (*JupiterDAO) GetDataByUID(uid int) (data *model.JupiterData, err error) {
	err = DB.Where("uid = ?", uid).First(&data).Error
	if err != nil {
		return nil, err
	}
	return
}

func (*JupiterDAO) GetNewestCfbp() ([]*model.JupiterData, error) {
	var dataList []*model.JupiterData
	err := DB.Order("cfbp_updated_at DESC").Find(&dataList).Error
	return dataList, err
}

func (*JupiterDAO) GetCourseByID(id int) (course *model.Course, err error) {
	err = DB.First(&course, id).Error
	return
}

func (*JupiterDAO) GetCoursesByUID(uid int) (courses []*model.Course, err error) {
	err = DB.Find(&courses, "uid = ?", uid).Error
	return
}

func (*JupiterDAO) GetAssignmentByID(id int) (assignment *model.Assignment, err error) {
	err = DB.First(&assignment, id).Error
	return
}

func (*JupiterDAO) GetAssignmentByInfo(title string, due *time.Time, from string) (assignment *model.Assignment, err error) {
	err = DB.Order("desc_fetched_at DESC").First(&assignment, "title = ? AND due = ? AND `from` = ?", title, due, from).Error
	return
}

func (*JupiterDAO) GetAssignmentsByUID(uid int) (assignments []*model.Assignment, err error) {
	err = DB.Order("due DESC").Find(&assignments, "uid = ?", uid).Error
	return
}

func (*JupiterDAO) GetAssignmentsByCourseAndUID(uid int, courseTitle string) (assignments []*model.Assignment, err error) {
	err = DB.Order("due DESC").Find(&assignments, "uid = ? AND `from` = ?", uid, courseTitle).Error
	return
}

func (*JupiterDAO) InsertData(data *model.JupiterData) error {
	return DB.Create(data).Error
}

func (*JupiterDAO) InsertCourse(course *model.Course) error {
	return DB.Create(course).Error
}

func (dao *JupiterDAO) InsertAssignment(old *model.Assignment, assignment *model.Assignment) error {
	// Check whether the assignment only changed date
	storedAssignments, _ := dao.GetAssignmentsByCourseAndUID(assignment.UID, assignment.From)
	for _, a := range storedAssignments {
		if a.Title == assignment.Title && a.Desc == assignment.Desc && a.Score == assignment.Score && old != nil {
			dao.UpdateAssignment(old, assignment)
			return nil
		}
	}

	return DB.Create(assignment).Error
}

func (*JupiterDAO) UpdateCourse(old *model.Course, course *model.Course) error {
	// Select * to select all columns, because status update can use 0 {Default will not update for gorm}
	return DB.Model(old).Select("*").Updates(course).Error
}

func (*JupiterDAO) UpdateAssignmentStatus(id int, status int) error {
	return DB.Model(&model.Assignment{ID: id}).Update("status", status).Error
}

func (*JupiterDAO) UpdateAssignmentNotFound(id int, notFound int) error {
	return DB.Model(&model.Assignment{ID: id}).Update("not_found", notFound).Error
}

func (*JupiterDAO) UpdateAssignment(old *model.Assignment, assignment *model.Assignment) error {
	// Select * to select all columns, because status update can use 0 {Default will not update for gorm}
	return DB.Model(old).Select("*").Updates(assignment).Error
}

func (*JupiterDAO) UpdateFetchTimeAndGPA(uid int, gpa string) error {
	if gpa == "" {
		return DB.Model(&model.JupiterData{UID: uid}).Where("uid = ?", uid).Update("fetched_at", time.Now()).Error
	}
	return DB.Model(&model.JupiterData{UID: uid}).Where("uid = ?", uid).Updates(map[string]any{"gpa": gpa, "fetched_at": time.Now()}).Error
}

func (*JupiterDAO) DeleteAssignment(id int) error {
	return DB.Delete(&model.Assignment{ID: id}).Error
}

package service

import (
	"adorable-star/internal/dao"
	"adorable-star/internal/model"
	"adorable-star/internal/pkg/crawler"
	"adorable-star/internal/pkg/response"
	"adorable-star/internal/pkg/util"
	"errors"
	"os"
	"strconv"
)

var Data = &DataService{dao.Jupiter}

type DataService struct {
	d *dao.JupiterDAO
}

func (s *DataService) FetchData(uid int) {
	go crawler.CrawlerJob(uid)
}

func (s *DataService) FetchAssignmentDesc(uid int, id int) error {
	// Get assignment data
	storedAssignment, err := s.d.GetAssignmentByID(id)
	if err != nil {
		return err
	}

	// Start fetching assignment description
	go func() {
		assignment := crawler.FetchAssignmentDesc(uid, storedAssignment)
		crawler.StoreAssignmentsData(uid, assignment.From, []*model.Assignment{assignment})
	}()

	return nil
}

func (s *DataService) GetCourses(uid int) (courses []*model.Course, err error) {
	courses, err = s.d.GetCoursesByUID(uid)
	return
}

func (s *DataService) GetCourse(uid int, id int) (course *response.Course, err error) {
	rawCourse, err := s.d.GetCourseByID(id)
	if err != nil {
		return
	}

	assignments, err := s.d.GetAssignmentsByCourseAndUID(uid, rawCourse.Title)
	if err != nil {
		return
	}

	course = &response.Course{
		Course:      *rawCourse,
		Assignments: assignments,
	}

	return
}

func (s *DataService) GetAssignments(uid int) (assignments []*model.Assignment, err error) {
	assignments, err = s.d.GetAssignmentsByUID(uid)
	return
}

func (s *DataService) GetAssignment(id int) (assignment *model.Assignment, err error) {
	assignment, err = s.d.GetAssignmentByID(id)
	return
}

func (s *DataService) GetReport(uid int) (file []byte, err error) {
	file, err = os.ReadFile(util.GetCwd() + "/storage/img/report/" + strconv.Itoa(uid) + ".png")
	if os.IsNotExist(err) {
		err = errors.New("fileNotExist")
	}
	return
}

func (s *DataService) GetMessages(uid int) (messages []*model.Message, err error) {
	messages, err = dao.Message.GetListByUID(uid)
	return
}

func (s *DataService) GetMessage(id int) (message *model.Message, err error) {
	message, err = dao.Message.GetByID(id)
	return
}

func (s *DataService) DeleteMessage(id int) error {
	return dao.Message.Delete(id)
}

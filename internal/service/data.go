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
	crawler.CrawlerJob(uid)
}

func (s *DataService) FetchAssignmentDetail(uid int, id int, force bool) error {
	// TODO: Limit request count for each user

	// Get assignment data
	storedAssignment, err := s.d.GetAssignmentByID(id)
	if err != nil {
		return err
	}

	// Start fetching assignment description
	var assignment *model.Assignment
	if force {
		assignment = crawler.FetchAssignmentDetail(uid, storedAssignment, true)
	} else {
		assignment = crawler.FetchAssignmentDetail(uid, storedAssignment)
	}

	crawler.StoreAssignmentsData(uid, assignment.From, []*model.Assignment{assignment}, true)

	return nil
}

func (s *DataService) GetData(uid int) (data *model.JupiterData, err error) {
	data, err = s.d.GetDataByUID(uid)
	return
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
	file, err = os.ReadFile(util.GetCwd() + "/storage/" + strconv.Itoa(uid) + "/report.png")
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

func (s *DataService) UpdateAssignmentStatus(id int, status int) error {
	return dao.Jupiter.UpdateAssignmentStatus(id, status)
}

func (s *DataService) TurnInJunoDoc(uid int, id int, text string) error {
	return crawler.TurnIn(uid, id, "JunoDoc", text)
}

func (s *DataService) TurnInFiles(uid int, id int) error {
	// Read user uploaded files
	entries, err := os.ReadDir(util.GetCwd() + "/storage/tmp" + strconv.Itoa(uid))
	if err != nil {
		return err
	}

	// Getting the file paths
	var files []string
	for _, entry := range entries {
		files = append(files, util.GetCwd()+"/storage/tmp"+strconv.Itoa(uid)+"/"+entry.Name())
	}

	// Turn in files
	err = crawler.TurnIn(uid, id, "File", files...)

	// Delete temporary stored files
	os.RemoveAll(util.GetCwd() + "/storage/tmp" + strconv.Itoa(uid))

	return err
}

func (s *DataService) UnSubmit(uid int, id int, name string) error {
	return crawler.UnSubmit(uid, id, name)
}

func (s *DataService) DeleteAllMessages(uid int) error {
	return dao.Message.DeleteAll(uid)
}

func (s *DataService) DeleteMessage(id int) error {
	return dao.Message.Delete(id)
}

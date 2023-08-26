package task

import (
	"cli/entity"
	"fmt"
)

type ServiceRepository interface {
	//DoesThisUserHaveThisCategoryID(userID, categoryID int) bool
	CreateNewTask(t entity.Task) (entity.Task, error)
	ListUserTask(userID int) ([]entity.Task, error)

	//domain driven

}
type Service struct {
	repository ServiceRepository
}

func NewService(repo ServiceRepository) *Service {
	return &Service{repository: repo}
}

type CreateRequest struct {
	Title               string
	DueDate             string
	CategoryID          int
	AuthenticatedUserID int
}
type CreateResponse struct {
	Task entity.Task
}

func (t *Service) Create(req CreateRequest) (CreateResponse, error) {

	//	if t.repository.DoesThisUserHaveThisCategoryID(req.AuthenticatedUserID, req.CategoryID) {
	//	return CreateResponse{}, fmt.Errorf("user doses not have this category:%d", req.CategoryID)
	//}

	createdTask, cErr := t.repository.CreateNewTask(entity.Task{
		ID:         0,
		Title:      req.Title,
		DueDate:    req.DueDate,
		CategoryID: req.CategoryID,
		IsDone:     false,
		UserID:     req.AuthenticatedUserID,
	})
	if cErr != nil {
		return CreateResponse{}, fmt.Errorf("cant't create New task: %v ", cErr)
	}
	return CreateResponse{Task: createdTask}, nil
}

type ListResponse struct {
	Tasks []entity.Task
}
type ListRequest struct {
	UserID int
}

func (t *Service) List(req ListRequest) (ListResponse, error) {
	taskList, err := t.repository.ListUserTask(req.UserID)
	if err != nil {
		return ListResponse{}, fmt.Errorf("can't list tasks: %v", err)
	}
	return ListResponse{Tasks: taskList}, nil

}

package api

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskAPI interface {
	AddTask(c *gin.Context)
	UpdateTask(c *gin.Context)
	DeleteTask(c *gin.Context)
	GetTaskByID(c *gin.Context)
	GetTaskList(c *gin.Context)
	GetTaskListByCategory(c *gin.Context)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskRepo service.TaskService) *taskAPI {
	return &taskAPI{taskRepo}
}

func (t *taskAPI) AddTask(c *gin.Context) {
	var newTask model.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	err := t.taskService.Store(&newTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "add task success"})
}

func (t *taskAPI) UpdateTask(c *gin.Context) {
	_, err := c.Request.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("unauthorized"))
		return
	}

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid task ID"})
		return
	}

	myTask, err := t.taskService.GetByID(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var taskAfterUpdate model.Task
	if err := c.ShouldBindJSON(&taskAfterUpdate); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	myTask.ID = taskAfterUpdate.ID
	myTask.Deadline = taskAfterUpdate.Deadline
	myTask.Priority = taskAfterUpdate.Priority
	myTask.CategoryID = taskAfterUpdate.CategoryID
	myTask.Status = taskAfterUpdate.Status

	err = t.taskService.Update(myTask.ID, myTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "update task success"})
}

func (t *taskAPI) DeleteTask(c *gin.Context) {
	_, err := c.Request.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("unauthorized"))
		return
	}

	idTask, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid task ID"})
		return
	}

	err = t.taskService.Delete(idTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "delete task success"})
}

func (t *taskAPI) GetTaskByID(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid task ID"})
		return
	}

	task, err := t.taskService.GetByID(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (t *taskAPI) GetTaskList(c *gin.Context) {
	_, err := c.Request.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("unauthorized"))
		return
	}

	tasksList, err := t.taskService.GetList()
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasksList)
}

func (t *taskAPI) GetTaskListByCategory(c *gin.Context) {
	_, err := c.Request.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("unauthorized"))
		return
	}

	idCat, _ := strconv.Atoi(c.Param("id"))

	taskListByCategory, err := t.taskService.GetTaskCategory(idCat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, taskListByCategory)
}

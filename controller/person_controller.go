package controller

import (
	"errors"
	"fmt"
	"infilon_task/models"
	"infilon_task/service"
	"infilon_task/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PersonController has dependency Injection for Service
type PersonController struct {
	Service *service.PersonService
}

// NewPersonController creates a new PersonController with the given service
func NewPersonController(service *service.PersonService) *PersonController {
	return &PersonController{
		Service: service,
	}
}

// GetPersonInfo handles the request to retrieve person information
func (pc *PersonController) GetPersonInfo(c *gin.Context) {
	personIDStr := c.Param("person_id")
	personID, err := strconv.Atoi(personIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid person ID: enter person id in integer"})
		return
	}

	// Call the service
	person, err := pc.Service.GetPersonInfo(personID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	if person.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Person not found"})
		return
	}

	c.JSON(http.StatusOK, person)
}

// CreatePerson handles the request to create a new person
func (pc *PersonController) CreatePerson(c *gin.Context) {
	var personRequest models.PersonRequest

	// Bind JSON body to the personRequest struct
	if err := c.BindJSON(&personRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}
	var person models.PersonRequest
	// Call the service
	person, err := pc.Service.CreatePerson(personRequest)
	if err != nil {
		if errors.Is(err, utils.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"message": fmt.Sprintf("Conflict error: %v", err)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Internal server error: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Person created successfully", "data": person})
}

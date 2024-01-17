package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"test-task/person/entity"
	"test-task/person/service"
)

type Handler struct {
	service.Service
}

func NewHandler(s service.Service) *Handler {
	return &Handler{s}
}

func (h *Handler) AddPerson(c *gin.Context) {
	var req *entity.PersonReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.AddPerson(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "person successfully added"})
}

func (h *Handler) GetPeople(c *gin.Context) {
	nation := c.Query("nationality")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))

	people, err := h.Service.GetPeople(nation, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, people)
}

func (h *Handler) DeletePerson(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.Service.DeletePerson(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "person delete successfully"})
}

func (h *Handler) GetPersonById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	person, err := h.Service.GetPersonById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, person)
}

func (h *Handler) UpdatePerson(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	updated := entity.Person{}

	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.Service.UpdatePerson(id, updated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error updating person: %s", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Person updated successfully"})
}

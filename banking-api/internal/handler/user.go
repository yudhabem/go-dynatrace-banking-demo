package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Random(c *gin.Context) {

	customer, account, err := h.service.CreateRandomUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"customerId":    customer.CustomerID,
			"name":          customer.Name,
			"accountNumber": account.AccountNumber,
			"balance":       account.Balance,
		},
	})
}

func (h *UserHandler) GetAll(c *gin.Context) {

	users, err := h.service.GetAllUsers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

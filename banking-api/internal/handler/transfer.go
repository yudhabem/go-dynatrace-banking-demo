package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/dto"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/service"
)

type TransferHandler struct {
	service *service.TransferService
}

func NewTransferHandler(service *service.TransferService) *TransferHandler {
	return &TransferHandler{
		service: service,
	}
}

func (h *TransferHandler) Transfer(c *gin.Context) {

	var req dto.TransferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	transactionID, err := h.service.Transfer(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"transactionId": transactionID,
	})
}

func (h *TransferHandler) Inquiry(c *gin.Context) {

	account := c.Param("account")

	acc, err := h.service.Inquiry(account)

	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    acc,
	})
}

func (h *TransferHandler) History(c *gin.Context) {

	data, err := h.service.History()

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

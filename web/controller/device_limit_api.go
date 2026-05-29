package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhsanaei/3x-ui/v3/database"
	"github.com/mhsanaei/3x-ui/v3/database/model"
	"github.com/mhsanaei/3x-ui/v3/logger"
	"github.com/mhsanaei/3x-ui/v3/web/service"
)

// GetClientDevices returns all devices for a client email
func (a *APIController) GetClientDevices(ctx *gin.Context) {
	email := ctx.Param("email")
	if email == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "email is required"})
		return
	}

	deviceSvc := &service.DeviceLimitService{}
	devices, err := deviceSvc.GetClientDevices(email)
	if err != nil {
		logger.Error("GetClientDevices error:", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":     "success",
		"devices": devices,
	})
}

// ClearClientDevice removes a specific device
func (a *APIController) ClearClientDevice(ctx *gin.Context) {
	email := ctx.Param("email")
	deviceID := ctx.Param("deviceId")

	if email == "" || deviceID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "email and deviceId are required"})
		return
	}

	deviceSvc := &service.DeviceLimitService{}
	err := deviceSvc.RemoveDevice(email, deviceID)
	if err != nil {
		logger.Error("ClearClientDevice error:", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "device removed"})
}

// ClearAllClientDevices removes all devices for a client
func (a *APIController) ClearAllClientDevices(ctx *gin.Context) {
	email := ctx.Param("email")
	if email == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "email is required"})
		return
	}

	db := database.GetDB()
	if err := db.Where("client_email = ?", email).Delete(&model.InboundClientDevices{}).Error; err != nil {
		logger.Error("ClearAllClientDevices error:", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "all devices cleared"})
}

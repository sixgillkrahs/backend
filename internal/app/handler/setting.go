package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetSetting(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "hello",
		})
	}
}

package routes

import (
	"github.com/gin-gonic/gin"
	"go-mnemosyne-api/category"
	"go-mnemosyne-api/note"
)

func UserRoute(routerGroup *gin.RouterGroup,
	categoryController category.Controller,
	noteController note.Controller) {
	categoryGroup := routerGroup.Group("/categories")
	categoryGroup.GET("", categoryController.GetAllByUser)
	categoryGroup.POST("", categoryController.Create)
	categoryGroup.PUT("", categoryController.Update)
	categoryGroup.DELETE(":id", categoryController.Delete)

	noteGroup := routerGroup.Group("/notes")
	noteGroup.POST("", noteController.Create)
	noteGroup.PUT(":id", noteController.Update)
	noteGroup.DELETE(":id", noteController.Delete)
}

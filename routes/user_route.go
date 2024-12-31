package routes

import (
	"github.com/gin-gonic/gin"
	"go-mnemosyne-api/category"
	"go-mnemosyne-api/note"
	"go-mnemosyne-api/user"
)

func UserRoute(routerGroup *gin.RouterGroup,
	userController user.Controller,
	categoryController category.Controller,
	noteController note.Controller) {
	categoryGroup := routerGroup.Group("/categories")
	categoryGroup.GET(":id", categoryController.GetAllByUser)
	categoryGroup.POST("", categoryController.Create)
	categoryGroup.PUT("", categoryController.Update)
	categoryGroup.DELETE(":id", categoryController.Delete)

	noteGroup := routerGroup.Group("/notes")
	noteGroup.POST("", noteController.Create)
}

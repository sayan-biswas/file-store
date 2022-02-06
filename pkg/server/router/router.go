package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sayan-biswas/file-store/pkg/database"
	"github.com/sayan-biswas/file-store/pkg/server/handler"
)

func Root(router *gin.Engine) {
	router.GET("/", handler.Root)
}

func Store(router *gin.Engine, store database.Store) {
	router.DELETE("/store", handler.RemoveFile(store))
	router.GET("/store", handler.GetFile(store))
	router.POST("/store", handler.AddFile(store))
	router.PUT("/store", handler.UpdateFile(store))
	router.GET("/store/check/file", handler.CheckFile(store))
	router.GET("/store/check/sha", handler.CheckSHA(store))
	router.GET("/store/list", handler.ListFiles(store))
	router.GET("/store/count", handler.WordCount(store))
	router.GET("/store/frequency", handler.WordFrequency(store))
}

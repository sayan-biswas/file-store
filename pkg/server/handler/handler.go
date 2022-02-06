package handler

import (
	_ "embed"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"github.com/sayan-biswas/file-store/pkg/database"
	"github.com/sayan-biswas/file-store/pkg/server/logger"
)

//go:embed store.info
var storeInfo string

type Frequency struct {
	Word  string
	Count int64
}

var log = logger.Log

const fingerprint string = "703273357638792F"

func Root(ctx *gin.Context) {
	ctx.Header("Store", fingerprint)
	ctx.String(http.StatusOK, storeInfo)
}

func GetFile(store database.Store) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		file := ctx.Query("file")
		data, err := store.Get(file)
		if err != nil {
			ctx.AbortWithError(http.StatusNotFound, err)
			return
		}
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file))
		ctx.Data(http.StatusOK, "application/octet-stream", data)
	}
	return gin.HandlerFunc(fn)
}

func AddFile(store database.Store) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		SHA := []byte(ctx.Request.FormValue("SHA"))
		fileHeader, err := ctx.FormFile("file")
		if err != nil {
			log.Error(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			log.Error(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil {
			log.Error(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		switch store.Add(fileHeader.Filename, SHA, data) {
		case nil:
			ctx.Status(http.StatusCreated)
		case badger.ErrConflict:
			ctx.AbortWithStatus(http.StatusConflict)
		default:
			log.Error(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
	return gin.HandlerFunc(fn)
}

func UpdateFile(store database.Store) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		SHA := []byte(ctx.Request.FormValue("SHA"))
		fileHeader, err := ctx.FormFile("file")
		if err != nil {
			log.Error(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			log.Error(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil {
			log.Error(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		fileExists := store.FileExists(fileHeader.Filename)
		switch store.Update(fileHeader.Filename, SHA, data) {
		case nil:
			if fileExists {
				ctx.Status(http.StatusOK)
			} else {
				ctx.Status(http.StatusCreated)
			}
		default:
			log.Error(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
	return gin.HandlerFunc(fn)
}

func RemoveFile(store database.Store) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		file := ctx.Query("file")
		err := store.Remove(file)
		switch err {
		case badger.ErrKeyNotFound:
			ctx.AbortWithStatus(http.StatusNotFound)
		case nil:
			ctx.Status(http.StatusNoContent)
		default:
			log.Error(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
	return gin.HandlerFunc(fn)
}

func ListFiles(store database.Store) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		details := false
		if ctx.Request.URL.Query().Has("details") {
			details, _ = strconv.ParseBool(ctx.Query("details"))
		}
		list, err := store.List(details)
		if err != nil {
			log.Error(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, list)
	}
	return gin.HandlerFunc(fn)
}

func WordCount(store database.Store) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		count, err := store.WordCount()
		if err != nil {
			log.Error(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.String(http.StatusOK, strconv.FormatInt(count, 10))
	}
	return gin.HandlerFunc(fn)
}

func WordFrequency(store database.Store) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		order := ctx.Query("order")
		frequency, err := store.WordFrequency()
		if err != nil {
			log.Error(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		keys := make([]string, 0, len(frequency))
		for key := range frequency {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool {
			if order == "dsc" {
				return frequency[keys[i]] > frequency[keys[j]]
			}
			return frequency[keys[i]] < frequency[keys[j]]
		})
		limit, err := strconv.Atoi(ctx.Query("limit"))
		if err != nil {
			err = nil
			limit = len(keys)
		}
		response := []Frequency{}
		for index, key := range keys {
			if index == limit {
				break
			}
			response = append(response, Frequency{Word: key, Count: frequency[key]})
		}
		ctx.JSON(http.StatusOK, response)
	}
	return gin.HandlerFunc(fn)
}

func CheckSHA(store database.Store) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		SHA, err := hex.DecodeString(ctx.Query("sha"))
		if err != nil {
			log.Error(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if store.SHAExists(SHA) {
			ctx.Status(http.StatusOK)
			return
		}
		ctx.AbortWithStatus(http.StatusNotFound)
	}
	return gin.HandlerFunc(fn)
}

func CheckFile(store database.Store) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		defer ctx.Request.Body.Close()
		file := ctx.Query("file")
		if store.FileExists(file) {
			ctx.Status(http.StatusOK)
			return
		}
		ctx.AbortWithStatus(http.StatusNotFound)
	}
	return gin.HandlerFunc(fn)
}

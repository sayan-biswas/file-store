package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sayan-biswas/file-store/pkg/database"
	"github.com/sayan-biswas/file-store/pkg/server/router"
	"github.com/stretchr/testify/assert"
)

type Frequency struct {
	Word  string
	Count int64
}

func init() {
	cleanUp()
}

func cleanUp() {
	err := os.RemoveAll("temp")
	if err != nil {
		log.Fatal(err)
	}
}
func TestServer(t *testing.T) {
	defer cleanUp()

	db, err := database.New(&database.Config{Path: "temp"})
	if err != nil {
		t.Fatal(err)
	}
	var store database.Store = db
	defer store.Close()

	gin.SetMode(gin.TestMode)
	server := gin.Default()
	router.Store(server, store)

	// Test Data
	fileName := "testfile.txt"
	data := []byte("this is test data for testing server")
	SHA256 := sha256.Sum256(data)
	SHA := SHA256[:]
	wordCount := "7"
	list := []database.File{{
		Name:      fileName,
		SHA:       hex.EncodeToString(SHA),
		Size:      36,
		WordCount: 7,
	}}
	frequency := map[string]int64{
		"this":    1,
		"is":      1,
		"test":    1,
		"data":    1,
		"for":     1,
		"testing": 1,
		"server":  1,
	}

	t.Run("Add_File", func(t *testing.T) {
		var byteReader *bytes.Reader
		var ioWriter io.Writer
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		defer writer.Close()
		byteReader = bytes.NewReader(SHA)
		ioWriter, _ = writer.CreateFormField("SHA")
		if _, err := io.Copy(ioWriter, byteReader); err != nil {
			t.Fatal(err)
		}
		byteReader = bytes.NewReader(data)
		ioWriter, _ = writer.CreateFormFile("file", fileName)
		if _, err := io.Copy(ioWriter, byteReader); err != nil {
			t.Fatal(err)
		}
		writer.Close()
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/store", bytes.NewReader(body.Bytes()))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		server.ServeHTTP(rr, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rr.Code)
		newData, err := store.Get(fileName)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, data, newData)
	})

	t.Run("Check_File", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/store/check/file?file="+fileName, nil)
		server.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Check_SHA", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/store/check/sha?sha="+hex.EncodeToString(SHA), nil)
		server.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Get_File", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/store?file="+fileName, nil)
		server.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, data, rr.Body.Bytes())
	})

	t.Run("Count_Word", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/store/count", nil)
		server.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, wordCount, rr.Body.String())
	})

	t.Run("Word_Frequency", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/store/frequency", nil)
		server.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		res := []Frequency{}
		err := json.Unmarshal(rr.Body.Bytes(), &res)
		assert.NoError(t, err)
		freqMap := map[string]int64{}
		for _, item := range res {
			freqMap[item.Word] = item.Count
		}
		assert.Equal(t, frequency, freqMap)
	})

	t.Run("List_File", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/store/list?details=true", nil)
		server.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		response := []database.File{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, list, response)
	})

	t.Run("Update_File", func(t *testing.T) {
		data := []byte("this is updated test data for testing server")
		SHA256 := sha256.Sum256(data)
		SHA := SHA256[:]

		var byteReader *bytes.Reader
		var ioWriter io.Writer
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		defer writer.Close()
		byteReader = bytes.NewReader(SHA)
		ioWriter, _ = writer.CreateFormField("SHA")
		if _, err := io.Copy(ioWriter, byteReader); err != nil {
			t.Fatal(err)
		}
		byteReader = bytes.NewReader(data)
		ioWriter, _ = writer.CreateFormFile("file", fileName)
		if _, err := io.Copy(ioWriter, byteReader); err != nil {
			t.Fatal(err)
		}
		writer.Close()
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/store", bytes.NewReader(body.Bytes()))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		server.ServeHTTP(rr, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)

		updatedData, err := store.Get(fileName)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, data, updatedData)
	})

	t.Run("Remove_File", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/store?file="+fileName, nil)
		server.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNoContent, rr.Code)
		assert.Equal(t, false, store.FileExists(fileName))
	})

	t.Run("Remove_File_Not_Found", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/store?file=randomfile.txt", nil)
		server.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

}

package client

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"

	"github.com/spf13/cobra"
)

func init() {
	store.AddCommand(add)
}

var add = &cobra.Command{
	Use:   "add",
	Long:  "Add files to store",
	Short: "Add files",
	Args:  cobra.MinimumNArgs(1),
	Run:   addFile,
}

func addFile(cmd *cobra.Command, args []string) {
	wg := &sync.WaitGroup{}
	for _, arg := range args {
		wg.Add(1)
		go storeAdd(arg, wg)
	}
	wg.Wait()
}

func storeAdd(filePath string, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer file.Close()
	fileName := path.Base(file.Name())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fileExists, err := checkFile(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if fileExists {
		fmt.Fprintf(os.Stdout, "%s - File already present in store\n", fileName)
		return
	}

	SHA256 := sha256.Sum256(data)
	SHA := SHA256[:]
	SHAExists, err := checkSHA(SHA)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	var byteReader *bytes.Reader
	var ioWriter io.Writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()
	if SHAExists {
		data = nil
	} else {
		SHA = nil
	}
	byteReader = bytes.NewReader(SHA)
	ioWriter, _ = writer.CreateFormField("SHA")
	if _, err := io.Copy(ioWriter, byteReader); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	byteReader = bytes.NewReader(data)
	ioWriter, _ = writer.CreateFormFile("file", fileName)
	if _, err := io.Copy(ioWriter, byteReader); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	writer.Close()

	storeURL, err := url.Parse(storeConfig.URL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	storeURL.Path = "store"
	req, err := http.NewRequest(http.MethodPost, storeURL.String(), bytes.NewReader(body.Bytes()))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated:
		fmt.Fprintf(os.Stdout, "%s - Added successfully\n", fileName)
	default:
		fmt.Fprintf(os.Stdout, "%s - Failed to upload\n", fileName)
	}
}

package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var location string

var get = &cobra.Command{
	Use:   "get",
	Long:  "Get file from Store",
	Short: "Get File",
	Args:  cobra.ExactArgs(1),
	Run:   getFile,
}

func init() {
	store.AddCommand(get)
	get.Flags().StringVarP(&location, "path", "p", ".", "path = ./files")
}

func getFile(cmd *cobra.Command, args []string) {
	storeURL, err := url.Parse(storeConfig.URL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	storeURL.Path = "store"
	values := storeURL.Query()
	values.Add("file", args[0])
	storeURL.RawQuery = values.Encode()

	req, err := http.NewRequest(http.MethodGet, storeURL.String(), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if res.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stdout, res.Status)
		return
	}

	file, err := os.Create(path.Join(location, args[0]))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer file.Close()

	if _, err := io.Copy(file, res.Body); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer res.Body.Close()
}

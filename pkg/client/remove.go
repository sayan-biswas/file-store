package client

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	store.AddCommand(remove)
}

var remove = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Long:    "Remove files to store",
	Short:   "Remove files",
	Args:    cobra.ExactArgs(1),
	Run:     removeFile,
}

func removeFile(cmd *cobra.Command, args []string) {
	storeURL, err := url.Parse(storeConfig.URL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	storeURL.Path = "store"
	values := storeURL.Query()
	values.Add("file", args[0])
	storeURL.RawQuery = values.Encode()

	req, err := http.NewRequest(http.MethodDelete, storeURL.String(), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		fmt.Fprintln(os.Stdout, res.Status)
		return
	}
	fmt.Fprintln(os.Stdout, "File deleted successfully!")
}

package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var details bool

var list = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Long:    "List files from store",
	Short:   "List files",
	Args:    cobra.NoArgs,
	Run:     listFiles,
}

func init() {
	store.AddCommand(list)
	list.Flags().BoolVarP(&details, "details", "d", false, "details = true | false")
}

func listFiles(cmd *cobra.Command, args []string) {
	storeURL, err := url.Parse(storeConfig.URL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	storeURL.Path = "store/list"
	values := storeURL.Query()
	values.Add("details", strconv.FormatBool(details))
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
		fmt.Fprintln(os.Stdout, res.StatusCode)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	type File struct {
		Name      string
		SHA       string
		Size      int64
		WordCount int64
	}

	if details {
		var files []File
		if err := json.Unmarshal(body, &files); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Fprintf(os.Stdout, "%-20s %10s %10s \n", "FILE NAME", "BYTES", "WORDS")
		fmt.Fprintf(os.Stdout, "%-20s %10s %10s \n", "---------", "-----", "-----")
		for _, file := range files {
			fmt.Fprintf(os.Stdout, "%-20s %10d %10d \n", file.Name, file.Size, file.WordCount)
		}
	} else {
		var files []string
		if err := json.Unmarshal(body, &files); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Fprintf(os.Stdout, "%s \n", "FILE NAME")
		fmt.Fprintf(os.Stdout, "%s \n", "---------")
		for _, file := range files {
			fmt.Fprintf(os.Stdout, "%s \n", file)
		}
	}
}

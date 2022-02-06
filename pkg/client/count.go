package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	store.AddCommand(count)
}

var count = &cobra.Command{
	Use:     "count",
	Aliases: []string{"wc"},
	Long:    "Display total word count",
	Short:   "Word count",
	Args:    cobra.NoArgs,
	Run:     wordCount,
}

func wordCount(cmd *cobra.Command, args []string) {
	storeURL, err := url.Parse(storeConfig.URL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	storeURL.Path = "store/count"

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

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Fprintf(os.Stdout, "Word Count: %s\n", string(body))
}

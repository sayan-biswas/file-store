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

var frequency = &cobra.Command{
	Use:     "frequency",
	Aliases: []string{"freq-words"},
	Long:    "Print word frequrncy in the all the files",
	Short:   "Word frequency",
	Args:    cobra.NoArgs,
	Run:     wordFrequency,
}

var order string
var limit int

func init() {
	store.AddCommand(frequency)
	frequency.Flags().StringVarP(&order, "order", "o", "asc", "order = dsc | asc")
	frequency.Flags().IntVarP(&limit, "limit", "n", 10, "limit = 10")
}

func wordFrequency(cmd *cobra.Command, args []string) {
	storeURL, err := url.Parse(storeConfig.URL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	storeURL.Path = "store/frequency"
	values := storeURL.Query()
	values.Add("order", order)
	values.Add("limit", strconv.Itoa(limit))
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
		fmt.Fprintf(os.Stdout, res.Status)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	type Frequency struct {
		Word  string
		Count int64
	}
	var frequency []Frequency
	if err := json.Unmarshal(body, &frequency); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	for _, item := range frequency {
		fmt.Fprintf(os.Stdout, "%10d %s \n", item.Count, item.Word)
	}

}

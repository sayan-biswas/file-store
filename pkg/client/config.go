package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type Config struct {
	URL string
}

const fingerprint string = "703273357638792F"

var config = &cobra.Command{
	Use:   "config",
	Long:  "Configure store URL",
	Short: "Confiure store",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := setConfig(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	},
}

func init() {
	store.AddCommand(config)
}

func setConfig() error {
	var storeURL string
	fmt.Printf("Store URL: ")
	fmt.Scanln(&storeURL)
	res, err := client.Get(storeURL)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}
	if fingerprint != res.Header.Get("Store") {
		return errors.New("invalid Store URL")
	}
	configJSON, err := json.Marshal(&Config{URL: storeURL})
	if err != nil {
		return err
	}
	if err := os.WriteFile(".StoreConfig", configJSON, 0666); err != nil {
		return err
	}
	return nil
}

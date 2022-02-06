package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var store = &cobra.Command{
	Use: "store",
}

var storeConfig Config
var client = &http.Client{
	Timeout: time.Second * 10,
}

func Execute() {
	if _, err := os.Stat(".StoreConfig"); errors.Is(err, os.ErrNotExist) {
		if err := setConfig(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}
	configFile, err := os.ReadFile(".StoreConfig")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if err := json.Unmarshal(configFile, &storeConfig); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if err := store.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

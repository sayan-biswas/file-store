package client

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed version.txt
var versionInfo string

func init() {
	store.AddCommand(version)
	store.Version = versionInfo
}

var version = &cobra.Command{
	Use:   "version",
	Long:  "Display store version",
	Short: "Store version",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Store Version %s\n", versionInfo)
	},
}

package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "netcheck",
	Short: "Network diagnostic tool for DevOps engineers",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error :%v\n", err)
		os.Exit(1)
	}
}

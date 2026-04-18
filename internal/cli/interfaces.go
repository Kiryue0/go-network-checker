package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Kiryue0/go-network-checker/internal/network"
	"github.com/spf13/cobra"
)

var interfacesCmd = &cobra.Command{
	Use:   "interfaces",
	Short: "List all network interfaces",
	Run: func(cmd *cobra.Command, args []string) {

		interfaces, err := network.GetInterfaces()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tIP\tMAC\tSTATUS\tMTU")
		fmt.Fprintln(w, "----\t--\t---\t------\t---")

		for _, iface := range interfaces {
			status := "DOWN"
			if iface.IsUp {
				status = "UP"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n",
				iface.Name,
				iface.IPAddress,
				iface.MACAddress,
				status,
				iface.MTU,
			)
		}
		if err := w.Flush(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(interfacesCmd)
}

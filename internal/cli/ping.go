package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"text/tabwriter"

	"github.com/Kiryue0/go-network-checker/internal/network"
	"github.com/spf13/cobra"
)

var pingFlag int

var pingCmd = &cobra.Command{
	Use:   "ping [hosts]",
	Short: "Ping one or more hosts",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()
		pings := network.PingHosts(ctx, args, pingFlag)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "HOST\tSTATUS\tRTT\tPACKET LOSS\tDATE TIME")
		fmt.Fprintln(w, "----\t------\t---\t-----------\t--------")

		for _, ping := range pings {
			status := "DEAD"
			if ping.IsAlive {
				status = "ALIVE"
			}
			fmt.Fprintf(w, "%s\t%s\t%v\t%1.f%%\t%s\n",
				ping.Host,
				status,
				ping.RTT,
				ping.PacketLoss,
				ping.Timestamp.Format("2006-01-02 15:04:05"),
			)
		}
		w.Flush()
	},
}

func init() {
	pingCmd.Flags().IntVarP(&pingFlag, "count", "c", 5, "Number of ping packets to send")
	rootCmd.AddCommand(pingCmd)
}

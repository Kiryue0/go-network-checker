package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/Kiryue0/go-network-checker/internal/metrics"
	"github.com/Kiryue0/go-network-checker/internal/network"
	"github.com/spf13/cobra"
)

var portsFlag string
var timeout time.Duration

var scanCmd = &cobra.Command{
	Use:   "scan [hosts]",
	Short: "Scan ports on one or more hosts",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		metrics.StartServer(":2112")
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()
		ports, err := parsePorts(portsFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return
		}

		scans := network.ScanPorts(ctx, args, ports, timeout)

		for _, host := range args {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintf(w, "%s\tPORT\tSTATE\tSERVICE\tRESPONSE TIME\n", host)
			fmt.Fprintf(w, "%s\t----\t-----\t-------\t-------------\n", strings.Repeat("-", len(host)))
			for _, scan := range scans {
				if scan.Host != host {
					continue
				}
				state := "CLOSED"
				if scan.IsOpen {
					state = "OPEN"
				}
				fmt.Fprintf(w, "\t%d\t%s\t%s\t%v\n",
					scan.Port,
					state,
					scan.Service,
					scan.ResponseTime,
				)
			}
			w.Flush()
			fmt.Println()
		}
	},
}

func init() {
	scanCmd.Flags().StringVarP(&portsFlag, "ports", "p", "22,80,443,3306,5432,6379,8080,27017", "Ports to scan (e.g. 80 | 22,80,443 | 100-200)")
	scanCmd.Flags().DurationVarP(&timeout, "timeout", "t", 2*time.Second, "Connection timeout")
	rootCmd.AddCommand(scanCmd)
}

func parsePorts(portsFlag string) ([]int, error) {
	var result []int
	parts := strings.Split(portsFlag, ",")

	for _, part := range parts {
		if strings.Contains(part, "-") {
			bounds := strings.Split(part, "-")
			start, err := strconv.Atoi(bounds[0])
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", bounds[0])
			}
			end, err := strconv.Atoi(bounds[1])
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", bounds[1])
			}
			for p := start; p <= end; p++ {
				result = append(result, p)
			}
		} else {
			p, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", part)
			}
			result = append(result, p)
		}
	}
	return result, nil
}

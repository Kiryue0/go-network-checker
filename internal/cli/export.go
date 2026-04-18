package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Kiryue0/go-network-checker/internal/export"
	"github.com/Kiryue0/go-network-checker/internal/model"
	"github.com/Kiryue0/go-network-checker/internal/network"
	"github.com/spf13/cobra"
)

var outputDir string

var exportCmd = &cobra.Command{
	Use:   "export [hosts]",
	Short: "Scan ports and save results to JSON",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()
		ports, err := parsePorts(portsFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return
		}

		portReport := network.ScanPorts(ctx, args, ports, timeout)
		pings := network.PingHosts(ctx, args, 3)
		aliveHosts := 0
		for _, ping := range pings {
			if ping.IsAlive {
				aliveHosts++
			}
		}
		scanReport := model.ScanReport{
			ScanDate:    time.Now(),
			TotalHost:   len(args),
			PortResults: portReport,
			Results:     pings,
			AliveHost:   aliveHosts,
		}
		exportedError := export.SaveJSON(scanReport, outputDir)
		if exportedError != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", exportedError)
			return
		}

	},
}

func init() {
	exportCmd.Flags().StringVarP(&portsFlag, "ports", "p", "22,80,443,3306,5432,6379,8080,27017", "Ports to scan (e.g. 80 | 22,80,443 | 100-200)")
	exportCmd.Flags().DurationVarP(&timeout, "timeout", "t", 2*time.Second, "Connection timeout")
	exportCmd.Flags().StringVarP(&outputDir, "output", "o", "./output", "Output directory for JSON file")
	rootCmd.AddCommand(exportCmd)
}

package network

import (
	"context"
	"errors"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Kiryue0/go-network-checker/internal/model"
)

func PingHost(host string, count int) (model.PingResult, error) {

	isAlive := true
	var rtt time.Duration
	timeStamp := time.Now()
	var packetLoss float64
	cmd := exec.Command("ping", "-c", strconv.Itoa(count), host)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return model.PingResult{
			Host:       host,
			IsAlive:    false,
			PacketLoss: 100.0,
			Timestamp:  timeStamp,
		}, errors.New("ping failed")
	}
	lines := strings.Split(string(output), "\n")
	var pLossLine string
	var rttLine string

	for _, line := range lines {
		if strings.Contains(line, "packet loss") {
			pLossLine = line
		}
		if strings.Contains(line, "rtt") {
			rttLine = line
		}
	}

	if pLossLine == "" || rttLine == "" {
		return model.PingResult{
			Host:       host,
			IsAlive:    false,
			PacketLoss: 100.0,
			Timestamp:  timeStamp,
		}, errors.New("ping output format not recognized")
	}

	packetLoss, err = strconv.ParseFloat(strings.Trim(strings.Fields(pLossLine)[5], "%"), 64)
	if err != nil {
		slog.Warn("failed to parse packet loss", "host", host, "error", err)
	}
	rtt, err = time.ParseDuration(strings.Split(rttLine, "/")[4] + "ms")
	if err != nil {
		slog.Warn("failed to parse RTT", "host", host, "error", err)
	}

	result := model.PingResult{
		Host:       host,
		IsAlive:    isAlive,
		RTT:        rtt,
		PacketLoss: packetLoss,
		Timestamp:  timeStamp,
	}

	return result, nil

}

func PingHosts(ctx context.Context, hosts []string, count int) []model.PingResult {
	var wg sync.WaitGroup
	results := make(chan model.PingResult, len(hosts))
	for _, host := range hosts {
		wg.Add(1)
		go func(h string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}
			result, err := PingHost(h, count)
			if err != nil {
				slog.Warn("ping failed", "host", h, "error", err)
			}
			results <- result

		}(host)
	}

	wg.Wait()
	close(results)

	var pingResults []model.PingResult
	for result := range results {
		pingResults = append(pingResults, result)
	}
	return pingResults
}

package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/govdbot/govd/internal/config"
)

// GetBandwidthString returns formatted total network bandwidth usage and limit comparison if configured.
func GetBandwidthString() string {
	bytesSent, bytesRecv := getNetIOCounters()
	totalBandwidth := bytesSent + bytesRecv

	if config.Env.MonthlyBandwidth > 0 {
		mbw := uint64(config.Env.MonthlyBandwidth) * 1024 * 1024 * 1024
		pct := float64(totalBandwidth) / float64(mbw) * 100
		return fmt.Sprintf("%s / %s [%.1f%%]", FormatBytes(totalBandwidth), FormatBytes(mbw), pct)
	}

	return FormatBytes(totalBandwidth)
}

// FormatBytes formats byte counts into human readable strings.
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func getNetIOCounters() (bytesSent uint64, bytesRecv uint64) {
	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return 0, 0
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}
		if fields[0] == "Inter-|" || fields[0] == "face" {
			continue
		}
		// Interface name may end with colon, e.g. "eth0:" or "eth0: 12345"
		iface := fields[0]
		if strings.HasPrefix(iface, "lo") {
			continue
		}

		var rxStr, txStr string
		if strings.Contains(iface, ":") {
			parts := strings.Split(iface, ":")
			if parts[1] != "" {
				rxStr = parts[1]
				txStr = fields[8]
			} else {
				rxStr = fields[1]
				txStr = fields[9]
			}
		} else {
			rxStr = fields[1]
			txStr = fields[9]
		}

		rxBytes, errRx := strconv.ParseUint(rxStr, 10, 64)
		txBytes, errTx := strconv.ParseUint(txStr, 10, 64)
		if errRx == nil && errTx == nil {
			bytesRecv += rxBytes
			bytesSent += txBytes
		}
	}
	return bytesSent, bytesRecv
}

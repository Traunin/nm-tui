package nmcli

import (
	"os/exec"
	"strconv"
	"strings"
)

type WifiNet struct {
	SSID   string
	Signal int
}

func ScanWifi() ([]WifiNet, error) {
	out, err := exec.Command("nmcli", "-t", "-f", "SSID,SIGNAL", "dev", "wifi").Output()
	if err != nil {
		return nil, err
	}

	var results []WifiNet
	lines := strings.SplitSeq(string(out), "\n")
	for line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}

		signal, _ := strconv.Atoi(parts[1])
		results = append(results, WifiNet{
			SSID:   parts[0],
			Signal: signal,
		})
	}

	return results, nil
}

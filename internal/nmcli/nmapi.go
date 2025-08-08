// Package nmcli provides interaction with nmcli utility
package nmcli

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/alphameo/nm-tui/internal/logger"
)

const nm = "nmcli"

type WifiNet struct {
	SSID   string
	Signal int
}

// ScanWifi shows list of wifi-networks able to be connected
// CMD: nmcli -t -f SSID,SIGNAL dev wifi
func ScanWifi() ([]WifiNet, error) {
	out, err := exec.Command(nm, "-t", "-f", "SSID,SIGNAL", "dev", "wifi").Output()
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

// ConnectWifi connects to wifi-network with given ssid using given password.
// CMD: nmcli device wifi connect "<SSID>" password "<PASSWORD>"
func ConnectWifi(ssid *string, password *string) error {
	// INFO: after nmcli 1.48.10 connection via password not able with saved networks
	DeleteConnection(ssid)
	args := []string{"device", "wifi", "connect", *ssid, "password", *password}
	out, err := exec.Command(nm, args...).Output()
	logger.InfoLog.Println(nm, args, "\n", string(out), err)
	return err
}

// ConnectSaved connects to wifi-network with given ssid if its password is saved.
// CMD: nmcli connection up "<SSID>"
func ConnectSaved(ssid *string) error {
	args := []string{"connection", "up", *ssid}
	out, err := exec.Command(nm, args...).Output()
	logger.InfoLog.Println(nm, args, "\n", string(out), err)
	return err
}

// GetConnected gives table of saved connections.
// CMD: nmcli -t -f NAME connection show
func GetConnected() ([]string, error) {
	args := []string{"-t", "-f", "NAME", "connection", "show"}
	out, err := exec.Command(nm, args...).Output()
	if err != nil {
		return nil, err
	}
	result := strings.Split(string(out), "\n")
	logger.InfoLog.Println(nm, args, "\n", strings.Join(result, ", "), err)
	return result, nil
}

func CheckPassword(ssid *string) (string, error) {
	args := []string{"-s", "-g", "802-11-wireless-security.psk", "connection", "show", *ssid}
	out, err := exec.Command(nm, args...).Output()
	if err != nil {
		return "", err
	}
	logger.InfoLog.Println(nm, args, "\n", string(out), err)
	return string(out), nil
}

// DeleteConnection removes wifi-network with given ssid from saved connections.
// CMD: nmcli connection delete "<SSID>"
func DeleteConnection(ssid *string) error {
	args := []string{"connection", "delete", *ssid}
	out, err := exec.Command(nm, args...).Output()
	logger.InfoLog.Println(nm, args, "\n", string(out), err)
	return err
}

// ConnectVpn connects to vpn with given vpnName
// CMD: nmcli connection up id "<VPN_NAME>"
func ConnectVpn(vpnName *string) error {
	args := []string{"connection", "up", "id", *vpnName}
	out, err := exec.Command(nm, args...).Output()
	logger.InfoLog.Println(nm, args, "\n", string(out), err)
	return err
}

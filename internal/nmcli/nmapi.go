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

// WifiScan shows list of wifi-networks able to be connected
// CMD: nmcli -t -f SSID,SIGNAL dev wifi
func WifiScan() ([]WifiNet, error) {
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

// WifiConnect connects to wifi-network with given ssid using given password.
// CMD: nmcli device wifi connect "<SSID>" password "<PASSWORD>"
func WifiConnect(ssid *string, password *string) error {
	WifiDeleteConnection(ssid) // FIX: after nmcli 1.48.10 connection via password not able with saved networks
	args := []string{"device", "wifi", "connect", *ssid, "password", *password}
	out, err := exec.Command(nm, args...).Output()
	if err == nil {
		logger.InfoLog.Printf("Connected to wifi %s (%s %s): %s\n", *ssid, nm, args, string(out))
	} else {
		logger.ErrorLog.Printf("Error connecting to wifi %s (%s %s): %s\n", *ssid, nm, args, err.Error())
	}
	return err
}

// WifiConnectSaved connects to wifi-network with given ssid if its password is saved.
// CMD: nmcli connection up "<SSID>"
func WifiConnectSaved(ssid *string) error {
	args := []string{"connection", "up", *ssid}
	out, err := exec.Command(nm, args...).Output()
	if err == nil {
		logger.InfoLog.Printf("Connected to saved wifi %s (%s %s): %s\n", *ssid, nm, args, string(out))
	} else {
		logger.ErrorLog.Printf("Error connecting to saved wifi %s (%s %s): %s\n", *ssid, nm, args, err.Error())
	}
	return err
}

// WifiGetConnected gives table of saved connections.
// CMD: nmcli -t -f NAME connection show
func WifiGetConnected() ([]string, error) {
	args := []string{"-t", "-f", "NAME", "connection", "show"}
	out, err := exec.Command(nm, args...).Output()
	if err != nil {
		logger.ErrorLog.Printf("Error retreiving list of connected wifi-networks (%s %s): %s\n", nm, args, err.Error())
		return nil, err
	}
	result := strings.Split(string(out), "\n")
	logger.InfoLog.Printf("Got list of connetcted wifi-networks (%s %s)\n", nm, args)
	return result, nil
}

// WifiGetPassword gives password of saved wifi-network with given ssid
// CMD: nmcli -s -g 802-11-wireless-security.psk connection show "<SSID>"
func WifiGetPassword(ssid *string) (string, error) {
	args := []string{"-s", "-g", "802-11-wireless-security.psk", "connection", "show", *ssid}
	out, err := exec.Command(nm, args...).Output()
	if err != nil {
		logger.ErrorLog.Printf("Error retrieving password to wifi %s (%s %s): %s\n", *ssid, nm, args, err.Error())
		return "", err
	}
	pw := strings.Trim(string(out), " \n")
	logger.InfoLog.Printf("Got password to wifi %s (%s %s)\n", *ssid, nm, args)
	return pw, nil
}

// WifiDeleteConnection removes wifi-network with given ssid from saved connections.
// CMD: nmcli connection delete "<SSID>"
func WifiDeleteConnection(ssid *string) error {
	args := []string{"connection", "delete", *ssid}
	out, err := exec.Command(nm, args...).Output()
	logger.InfoLog.Println(nm, args, "\n", string(out), err)
	if err == nil {
		logger.InfoLog.Printf("Connection to wifi %s was deleted (%s %s): %s", *ssid, nm, args, string(out))
	} else {
		logger.ErrorLog.Printf("Error deleting connection to wifi %s (%s %s): %s\n", *ssid, nm, args, err.Error())
	}
	return err
}

// VpnConnect connects to vpn with given vpnName
// CMD: nmcli connection up id "<VPN_NAME>"
func VpnConnect(vpnName *string) error {
	args := []string{"connection", "up", "id", *vpnName}
	out, err := exec.Command(nm, args...).Output()
	if err == nil {
		logger.InfoLog.Printf("Connected to VPN %s (%s %s): %s\n", *vpnName, nm, args, string(out))
	} else {
		logger.ErrorLog.Printf("Error connecting to VPN %s (%s %s): %s\n", *vpnName, nm, args, err.Error())
	}
	return err
}

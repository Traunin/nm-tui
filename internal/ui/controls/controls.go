// Package controls provides simple public controls of main ui
package controls

import (
	"github.com/alphameo/nm-tui/internal/nmcli"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	PopupContentMsg  tea.Model
	PopupActivityMsg bool
)

func SetPopupContent(content tea.Model) tea.Cmd {
	return func() tea.Msg {
		return PopupContentMsg(content)
	}
}

func SetPopupActivity(isActive bool) tea.Cmd {
	return func() tea.Msg {
		return PopupActivityMsg(isActive)
	}
}

type (
	NotificationTextMsg     string
	NotificationActivityMsg bool
)

func SetNotificationText(text string) tea.Cmd {
	return func() tea.Msg {
		return NotificationTextMsg(text)
	}
}

func SetNotificationActivity(isActive bool) tea.Cmd {
	return func() tea.Msg {
		return NotificationActivityMsg(isActive)
	}
}

func DeleteConnection(ssid string) tea.Cmd {
	return func() tea.Msg {
		nmcli.WifiDeleteConnection(ssid)
		return nil
	}
}

package vault

import (
	tea "github.com/charmbracelet/bubbletea"
)

type createSecretCommand struct {
}

func newCreateSecretCommand() createSecretCommand {
	return createSecretCommand{}
}

func (c createSecretCommand) execute() tea.Msg {
	return createSecretRequestedMsg{}
}

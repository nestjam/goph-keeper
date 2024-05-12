package auth

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nestjam/goph-keeper/internal/tui/vault"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginModel_Init(t *testing.T) {
	sut := NewLoginModel()

	got := sut.Init()

	assertEqualCmd(t, textinput.Blink, got)
}

func TestLoginModel_Update(t *testing.T) {
	t.Run("user typed server address", func(t *testing.T) {
		sut := NewLoginModel()
		const input = "http://localhost:8080"
		msg := tea.KeyMsg{Runes: []rune(input)}

		model, _ := sut.Update(msg)

		got, ok := model.(loginModel)
		assert.True(t, ok)
		assert.Equal(t, input, got.textinput.Value())
	})
	t.Run("user exited by ctrl+c", func(t *testing.T) {
		sut := NewLoginModel()
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}

		_, cmd := sut.Update(msg)

		assertEqualCmd(t, tea.Quit, cmd)
	})
	t.Run("user exited by esc", func(t *testing.T) {
		sut := NewLoginModel()
		msg := tea.KeyMsg{Type: tea.KeyEsc}

		_, cmd := sut.Update(msg)

		assertEqualCmd(t, tea.Quit, cmd)
	})
	t.Run("user entered server address", func(t *testing.T) {
		sut := tea.Model(NewLoginModel())

		const input = "http://localhost:8080"
		msg := tea.KeyMsg{Runes: []rune(input)}
		sut, _ = sut.Update(msg)

		msg = tea.KeyMsg{Type: tea.KeyEnter}
		model, cmd := sut.Update(msg)

		got, _ := model.(loginModel)
		assert.Equal(t, input, got.address)
		assert.Equal(t, "", got.textinput.Value())
		assert.Equal(t, "Enter email", got.textinput.Placeholder)
		assert.Nil(t, cmd)
	})
	t.Run("user pressed enter with empty input", func(t *testing.T) {
		sut := NewLoginModel()
		msg := tea.KeyMsg{Type: tea.KeyEnter}

		model, cmd := sut.Update(msg)

		got, _ := model.(loginModel)
		assert.Equal(t, "", got.address)
		assert.Equal(t, "Enter server address", got.textinput.Placeholder)
		assert.Nil(t, cmd)
	})
	t.Run("user entered email", func(t *testing.T) {
		m := NewLoginModel()
		m.address = "localhost:8080"
		sut := tea.Model(m)

		const input = "user@email.com"
		msg := tea.KeyMsg{Runes: []rune(input)}
		sut, _ = sut.Update(msg)

		msg = tea.KeyMsg{Type: tea.KeyEnter}
		model, cmd := sut.Update(msg)

		got, _ := model.(loginModel)
		assert.Equal(t, input, got.email)
		assert.Equal(t, "", got.textinput.Value())
		assert.Equal(t, "Enter password", got.textinput.Placeholder)
		assert.Nil(t, cmd)
	})
	t.Run("user entered password", func(t *testing.T) {
		m := NewLoginModel()
		m.address = "localhost:8080"
		m.email = "user@mail.com"
		sut := tea.Model(m)

		const input = "1234"
		msg := tea.KeyMsg{Runes: []rune(input)}
		sut, _ = sut.Update(msg)

		msg = tea.KeyMsg{Type: tea.KeyEnter}
		model, cmd := sut.Update(msg)

		got, _ := model.(loginModel)
		assert.Equal(t, input, got.password)
		loginCmd := loginCommand{}
		assertEqualCmd(t, loginCmd.execute, cmd)
	})
	t.Run("login completed", func(t *testing.T) {
		m := NewLoginModel()
		sut := tea.Model(m)
		msg := loginCompletedMsg{}

		model, cmd := sut.Update(msg)

		_, ok := model.(vault.SecretsModel)
		assert.True(t, ok)
		listCommand := vault.NewListSecretsCommand("", nil)
		assertEqualCmd(t, listCommand.Execute, cmd)
	})
	t.Run("register completed", func(t *testing.T) {
		m := NewLoginModel()
		sut := tea.Model(m)
		msg := registerCompletedMsg{}

		model, cmd := sut.Update(msg)

		_, ok := model.(vault.SecretsModel)
		assert.True(t, ok)
		listCommand := vault.NewListSecretsCommand("", nil)
		assertEqualCmd(t, listCommand.Execute, cmd)
	})
	t.Run("error on login", func(t *testing.T) {
		sut := NewLoginModel()
		msg := errMsg{errors.New("error")}

		model, _ := sut.Update(msg)

		got, _ := model.(loginModel)
		assert.Equal(t, msg.err, got.err)
	})
	t.Run("failed to login", func(t *testing.T) {
		m := NewLoginModel()
		m.address = "localhost:8080"
		m.email = "user@mail.com"
		m.password = "1234"
		sut := tea.Model(m)
		msg := loginFailedMsg{statusCode: http.StatusUnauthorized}

		model, _ := sut.Update(msg)

		got, _ := model.(loginModel)
		assert.Empty(t, got.email)
		assert.Empty(t, got.password)
	})
	t.Run("failed to register", func(t *testing.T) {
		m := NewLoginModel()
		m.address = "localhost:8080"
		m.email = "user@mail.com"
		m.password = "1234"
		sut := tea.Model(m)
		msg := registerFailedMsg{statusCode: http.StatusUnauthorized}

		model, _ := sut.Update(msg)

		got, _ := model.(loginModel)
		assert.Empty(t, got.email)
		assert.Empty(t, got.password)
	})
	t.Run("key down pressed from login choice", func(t *testing.T) {
		msg := tea.KeyMsg{Type: tea.KeyDown}
		const want = 1
		sut := NewLoginModel()

		model, cmd := sut.Update(msg)

		m, ok := model.(loginModel)
		assert.True(t, ok)
		got := m.cursor
		assert.Equal(t, want, got)
		assert.Nil(t, cmd)
	})
	t.Run("key up pressed from login choice", func(t *testing.T) {
		msg := tea.KeyMsg{Type: tea.KeyUp}
		const want = 1
		sut := NewLoginModel()

		model, cmd := sut.Update(msg)

		m, ok := model.(loginModel)
		assert.True(t, ok)
		got := m.cursor
		assert.Equal(t, want, got)
		assert.Nil(t, cmd)
	})
	t.Run("key down pressed from register choice", func(t *testing.T) {
		msg := tea.KeyMsg{Type: tea.KeyDown}
		const want = 0
		sut := NewLoginModel()
		sut.cursor = 1

		model, cmd := sut.Update(msg)

		m, ok := model.(loginModel)
		assert.True(t, ok)
		got := m.cursor
		assert.Equal(t, want, got)
		assert.Nil(t, cmd)
	})
	t.Run("key up pressed from register choice", func(t *testing.T) {
		msg := tea.KeyMsg{Type: tea.KeyUp}
		const want = 0
		sut := NewLoginModel()
		sut.cursor = 1

		model, cmd := sut.Update(msg)

		m, ok := model.(loginModel)
		assert.True(t, ok)
		got := m.cursor
		assert.Equal(t, want, got)
		assert.Nil(t, cmd)
	})
	t.Run("user pressed enter to register", func(t *testing.T) {
		m := NewLoginModel()
		m.address = "localhost:8080"
		m.email = "user@mail.com"
		m.password = "1234"
		m.cursor = 1
		sut := tea.Model(m)
		msg := tea.KeyMsg{Type: tea.KeyEnter}

		model, cmd := sut.Update(msg)

		_, ok := model.(loginModel)
		assert.True(t, ok)
		registerCmd := registerCommand{}
		assertEqualCmd(t, registerCmd.execute, cmd)
	})
	t.Run("window size changed", func(t *testing.T) {
		sut := NewLoginModel()
		msg := tea.WindowSizeMsg{Width: 100}
		require.NotEqual(t, msg.Width, sut.help.Width)

		model, _ := sut.Update(msg)

		got, ok := model.(loginModel)
		assert.True(t, ok)
		assert.Equal(t, msg.Width, got.help.Width)
	})
}

func assertEqualCmd(t *testing.T, want, got tea.Cmd) {
	t.Helper()

	gotValue := reflect.ValueOf(got)
	wantValue := reflect.ValueOf(want)
	assert.Equal(t, gotValue.Pointer(), wantValue.Pointer())
}

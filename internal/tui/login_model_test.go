package tui

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, input, got.serverAddress)
		assert.Equal(t, "", got.textinput.Value())
		assert.Equal(t, "Enter email", got.textinput.Placeholder)
		assert.Nil(t, cmd)
	})
	t.Run("user pressed enter with empty input", func(t *testing.T) {
		sut := NewLoginModel()
		msg := tea.KeyMsg{Type: tea.KeyEnter}

		model, cmd := sut.Update(msg)

		got, _ := model.(loginModel)
		assert.Equal(t, "", got.serverAddress)
		assert.Equal(t, "Enter server address", got.textinput.Placeholder)
		assert.Nil(t, cmd)
	})
	t.Run("user entered email", func(t *testing.T) {
		m := NewLoginModel()
		m.serverAddress = "localhost:8080"
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
		m.serverAddress = "localhost:8080"
		m.email = "user@mail.com"
		sut := tea.Model(m)

		const input = "1234"
		msg := tea.KeyMsg{Runes: []rune(input)}
		sut, _ = sut.Update(msg)

		msg = tea.KeyMsg{Type: tea.KeyEnter}
		model, cmd := sut.Update(msg)

		got, _ := model.(loginModel)
		assert.Equal(t, input, got.password)
		assert.Nil(t, cmd)
	})
}

func assertEqualCmd(t *testing.T, want, got tea.Cmd) {
	t.Helper()

	gotValue := reflect.ValueOf(got)
	wantValue := reflect.ValueOf(want)
	assert.Equal(t, gotValue.Pointer(), wantValue.Pointer())
}

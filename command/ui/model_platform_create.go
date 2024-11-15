package ui

import (
	"buildenv/command"
	"buildenv/config"
	"buildenv/pkg/color"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func newPlatformCreateModel(callbacks config.PlatformCallbacks, goback func(this *platformCreateModel)) *platformCreateModel {
	ti := textinput.New()
	ti.Placeholder = "for example: x86_64-linux-ubuntu-20.04..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 100
	ti.TextStyle = styleImpl.focusedStyle
	ti.PromptStyle = styleImpl.focusedStyle
	ti.Cursor.Style = styleImpl.focusedStyle

	return &platformCreateModel{
		textInput: ti,
		callbacks: callbacks,
		goback:    goback,
	}
}

type platformCreateModel struct {
	textInput textinput.Model
	created   bool
	err       error

	callbacks config.PlatformCallbacks
	goback    func(this *platformCreateModel)
}

func (p platformCreateModel) Init() tea.Cmd {
	return textinput.Blink
}

func (p platformCreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return p, tea.Quit

		case "esc":
			p.goback(&p)
			return p, nil

		case "enter":
			if err := p.callbacks.OnCreatePlatform(p.textInput.Value()); err != nil {
				p.err = err
				p.created = false
			} else {
				p.created = true
			}

			return p, tea.Quit
		}
	}

	var cmd tea.Cmd
	p.textInput, cmd = p.textInput.Update(msg)
	return p, cmd
}

func (p platformCreateModel) View() string {
	if p.created {
		return command.PlatformCreated(p.textInput.Value())
	}

	if p.err != nil {
		return command.PlatformCreateFailed(p.textInput.Value(), p.err)
	}

	return fmt.Sprintf("\n%s\n\n%s\n\n%s\n",
		color.Sprintf(color.Blue, "Please input your platform name: "),
		p.textInput.View(),
		color.Sprintf(color.Gray, "[esc -> back | ctrl+c/q -> quit]"),
	)
}

func (p *platformCreateModel) Reset() {
	p.textInput.Reset()
	p.created = false
	p.err = nil
}

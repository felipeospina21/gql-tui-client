package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
)

// Available spinners
var spinners = []spinner.Spinner{
	spinner.Line,
	spinner.Dot,
	spinner.MiniDot,
	spinner.Jump,
	spinner.Pulse,
	spinner.Points,
	spinner.Globe,
	spinner.Moon,
	spinner.Monkey,
}

type spinnerModel struct {
	model spinner.Model
	index int
}

func (m *mainModel) newSpinnerModel() {
	m.spinner.model = spinner.New()
	m.spinner.model.Style = spinnerStyle
	// s.Spinner = spinner.Pulse
}

func (m *mainModel) nextSpinner() {
	if m.spinner.index == len(spinners)-1 {
		m.spinner.index = 0
	} else {
		m.spinner.index++
	}
}

func (m *mainModel) resetSpinner() {
	m.spinner.model = spinner.New()
	m.spinner.model.Style = spinnerStyle
	m.spinner.model.Spinner = spinners[m.spinner.index]
}

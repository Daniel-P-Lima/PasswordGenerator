package main

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type appState int

const (
	stLen appState = iota
	stDigits
	stSymbols
	stNoUpper
	stAllowRepeat
	stReview
	stGenerating
	stDone
)

type keymap struct {
	Next    key.Binding
	Prev    key.Binding
	Confirm key.Binding
	Back    key.Binding
	Toggle  key.Binding
	Quit    key.Binding
}

func newKeymap() keymap {
	return keymap{
		Next:    key.NewBinding(key.WithKeys("tab", "right"), key.WithHelp("tab", "próximo")),
		Prev:    key.NewBinding(key.WithKeys("shift+tab", "left"), key.WithHelp("S-tab", "anterior")),
		Confirm: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirmar")),
		Back:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "voltar")),
		Toggle:  key.NewBinding(key.WithKeys(" "), key.WithHelp("espaço", "alternar")),
		Quit:    key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("q", "sair")),
	}
}

type Model struct {
	state appState
	keys  keymap
	help  help.Model

	// inputs
	inLen       textinput.Model
	inDigits    textinput.Model
	inSymbols   textinput.Model
	noUpper     bool
	allowRepeat bool

	// validação
	err string

	// resultado
	result string
}

func NewModel() Model {
	mk := newKeymap()
	h := help.New()
	h.ShowAll = false

	mkIn := func(ph string, val string) textinput.Model {
		in := textinput.New()
		in.Prompt = "➜ "
		in.Placeholder = ph
		in.SetValue(val)
		in.CharLimit = 4
		in.Focus()
		return in
	}

	m := Model{
		state:       stLen,
		keys:        mk,
		help:        h,
		inLen:       mkIn("Tamanho (p.ex. 16)", "16"),
		inDigits:    mkIn("Qtd dígitos (0..)", "4"),
		inSymbols:   mkIn("Qtd símbolos (0..)", "2"),
		noUpper:     false,
		allowRepeat: false,
	}
	return m
}

func (m Model) Init() tea.Cmd { return textinput.Blink }

func (m *Model) validateIntField(in *textinput.Model, min int) (int, bool) {
	v := in.Value()
	if v == "" {
		m.err = "campo obrigatório"
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		m.err = "use apenas números"
		return 0, false
	}
	if n < min {
		m.err = fmt.Sprintf("deve ser ≥ %d", min)
		return 0, false
	}
	m.err = ""
	return n, true
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Next):
			return m.handleNext()

		case key.Matches(msg, m.keys.Prev):
			return m.handlePrev()

		case key.Matches(msg, m.keys.Back):
			return m.handleBack()

		case key.Matches(msg, m.keys.Toggle):
			return m.handleToggle()

		case key.Matches(msg, m.keys.Confirm):
			return m.handleConfirm()
		}
	}

	var cmd tea.Cmd
	switch m.state {
	case stLen:
		m.inLen, cmd = m.inLen.Update(msg)
	case stDigits:
		m.inDigits, cmd = m.inDigits.Update(msg)
	case stSymbols:
		m.inSymbols, cmd = m.inSymbols.Update(msg)
	}
	return m, cmd
}

func (m Model) handleNext() (tea.Model, tea.Cmd) {
	switch m.state {
	case stLen:
		return m.to(stDigits)
	case stDigits:
		return m.to(stSymbols)
	case stSymbols:
		return m.to(stNoUpper)
	case stNoUpper:
		return m.to(stAllowRepeat)
	case stAllowRepeat:
		return m.to(stReview)
	}
	return m, nil
}

func (m Model) handlePrev() (tea.Model, tea.Cmd) {
	switch m.state {
	case stDigits:
		return m.to(stLen)
	case stSymbols:
		return m.to(stDigits)
	case stNoUpper:
		return m.to(stSymbols)
	case stAllowRepeat:
		return m.to(stNoUpper)
	case stReview:
		return m.to(stAllowRepeat)
	case stDone:
		return m.to(stReview)
	}
	return m, nil
}

func (m Model) handleBack() (tea.Model, tea.Cmd) { return m.handlePrev() }

func (m Model) handleToggle() (tea.Model, tea.Cmd) {
	switch m.state {
	case stNoUpper:
		m.err = ""
		m.noUpper = !m.noUpper
	case stAllowRepeat:
		m.err = ""
		m.allowRepeat = !m.allowRepeat
	}
	return m, nil
}

func (m Model) handleConfirm() (tea.Model, tea.Cmd) {
	switch m.state {
	case stLen:
		if _, ok := m.validateIntField(&m.inLen, 1); ok {
			return m.handleNext()
		}
	case stDigits:
		if _, ok := m.validateIntField(&m.inDigits, 0); ok {
			return m.handleNext()
		}
	case stSymbols:
		if _, ok := m.validateIntField(&m.inSymbols, 0); ok {
			return m.handleNext()
		}
	case stNoUpper:
		return m.handleNext()
	case stAllowRepeat:
		return m.handleNext()
	case stReview:
		// gerar
		m.state = stGenerating
		cfg := GenConfig{}
		cfg.Length, _ = strconv.Atoi(m.inLen.Value())
		cfg.NumDigits, _ = strconv.Atoi(m.inDigits.Value())
		cfg.NumSymbols, _ = strconv.Atoi(m.inSymbols.Value())
		cfg.NoUpper = m.noUpper
		cfg.AllowRepeat = m.allowRepeat
		m.result = Generate(cfg)
		m.state = stDone
	case stDone:
		// regenerar com mesmas opções
		m.state = stGenerating
		cfg := GenConfig{}
		cfg.Length, _ = strconv.Atoi(m.inLen.Value())
		cfg.NumDigits, _ = strconv.Atoi(m.inDigits.Value())
		cfg.NumSymbols, _ = strconv.Atoi(m.inSymbols.Value())
		cfg.NoUpper = m.noUpper
		cfg.AllowRepeat = m.allowRepeat
		m.result = Generate(cfg)
		_ = SavePassword(m.result)
		m.state = stDone
	}
	return m, nil
}

func (m Model) to(s appState) (tea.Model, tea.Cmd) { m.state = s; m.err = ""; return m, nil }

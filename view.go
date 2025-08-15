package main

import (
	"fmt"
	"strings"
)

func checkbox(label string, v bool, focused bool) string {
	box := "[ ]"
	if v {
		box = "[x]"
	}
	line := fmt.Sprintf("%s %s", box, label)
	if focused {
		return focusBox.Render(line)
	}
	return boxStyle.Render(line)
}

func viewHeader() string {
	return titleStyle.Render("Password Generator") + "\n" + subtleStyle.Render("Use Tab/Shift+Tab para navegar, Enter para confirmar, \nQ para sair") + "\n\n"
}

func (m Model) viewStatus() string {
	steps := []string{"Len", "Digits", "Symbols", "NoUpper", "Repeat", "Review"}
	cur := 0
	switch m.state {
	case stLen:
		cur = 0
	case stDigits:
		cur = 1
	case stSymbols:
		cur = 2
	case stNoUpper:
		cur = 3
	case stAllowRepeat:
		cur = 4
	case stReview, stGenerating, stDone:
		cur = 5
	}
	for i := range steps {
		if i == cur {
			steps[i] = okStyle.Render(steps[i])
		} else {
			steps[i] = subtleStyle.Render(steps[i])
		}
	}
	return statusStyle.Render(" " + strings.Join(steps, " • ") + " ")
}

func (m Model) viewError() string {
	if m.err == "" {
		return ""
	}
	return "\n" + errorStyle.Render("✖ "+m.err) + "\n"
}

func (m Model) viewReview() string {
	return boxStyle.Render(fmt.Sprintf(
		"Tamanho: %s\nDígitos: %s\nSímbolos: %s\nSem maiúsculas: %v\nPermitir repetição: %v",
		m.inLen.Value(), m.inDigits.Value(), m.inSymbols.Value(), m.noUpper, m.allowRepeat,
	))
}

func (m Model) View() string {
	b := strings.Builder{}
	b.WriteString(viewHeader())

	switch m.state {
	case stLen:
		b.WriteString(labelStyle.Render("1/5 • Tamanho"))
		b.WriteString("\n" + focusBox.Render(m.inLen.View()))
		b.WriteString(m.viewError())

	case stDigits:
		b.WriteString(labelStyle.Render("2/5 • Dígitos"))
		b.WriteString("\n" + focusBox.Render(m.inDigits.View()))
		b.WriteString(m.viewError())

	case stSymbols:
		b.WriteString(labelStyle.Render("3/5 • Símbolos"))
		b.WriteString("\n" + focusBox.Render(m.inSymbols.View()))
		b.WriteString(m.viewError())

	case stNoUpper:
		b.WriteString(labelStyle.Render("4/5 • Sem maiúsculas?"))
		b.WriteString("\n" + checkbox("Não usar letras maiúsculas", m.noUpper, true))
		b.WriteString("\n" + subtleStyle.Render("Use espaço para alternar"))

	case stAllowRepeat:
		b.WriteString(labelStyle.Render("5/5 • Permitir repetição?"))
		b.WriteString("\n" + checkbox("Permitir caracteres repetidos", m.allowRepeat, true))
		b.WriteString("\n" + subtleStyle.Render("Use espaço para alternar"))

	case stReview:
		b.WriteString(labelStyle.Render("Revisar opções"))
		b.WriteString("\n" + m.viewReview())
		b.WriteString("\n\n" + subtleStyle.Render("Enter para gerar • Esc para voltar"))

	case stGenerating:
		b.WriteString(labelStyle.Render("Gerando..."))

	case stDone:
		b.WriteString(labelStyle.Render("Senha gerada"))
		b.WriteString("\n" + focusBox.Render(m.result))
		b.WriteString("\n\n" + subtleStyle.Render("Enter para regerar • Esc para alterar opções • q para sair"))
	}

	b.WriteString("\n\n" + m.viewStatus() + "\n")
	return b.String()
}

package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	cTitle  = lipgloss.Color("63")
	cSubtle = lipgloss.Color("245")
	cError  = lipgloss.Color("203")
	cOk     = lipgloss.Color("42")
	cFocus  = lipgloss.Color("105")
	cDimBg  = lipgloss.Color("236")

	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(cTitle)
	subtleStyle = lipgloss.NewStyle().Foreground(cSubtle)
	errorStyle  = lipgloss.NewStyle().Foreground(cError)
	okStyle     = lipgloss.NewStyle().Foreground(cOk)
	labelStyle  = lipgloss.NewStyle().Foreground(cFocus)
	statusStyle = lipgloss.NewStyle().Foreground(cSubtle).Background(cDimBg).Padding(0, 1)
	boxStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2)
	focusBox    = boxStyle.Copy().BorderForeground(cFocus)
)

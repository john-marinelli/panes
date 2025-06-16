package panes

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const borderWidth int = 1

type resizer func(msg tea.WindowSizeMsg) tea.WindowSizeMsg

type pos struct {
	x int
	y int
}

type Styles struct {
	ActiveBorder lipgloss.Border

	HorizontalJoinPosition lipgloss.Position
	VerticalJoinPosition   lipgloss.Position

	ActiveStyle   lipgloss.Style
	InactiveStyle lipgloss.Style
}

type KeyMap struct {
	Left  key.Binding
	Right key.Binding
	Down  key.Binding
	Up    key.Binding
	Quit  key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Left:  key.NewBinding(key.WithKeys("ctrl+h")),
		Right: key.NewBinding(key.WithKeys("ctrl+l")),
		Down:  key.NewBinding(key.WithKeys("ctrl+j")),
		Up:    key.NewBinding(key.WithKeys("ctrl+k")),
		Quit:  key.NewBinding(key.WithKeys("ctrl+c")),
	}
}

type Model struct {
	sections [][]tea.Model
	resizers []resizer
	active   pos

	KeyMap KeyMap
	Styles Styles
}

func New(sections [][]tea.Model) Model {
	rs := []resizer{}
	for _, row := range sections {
		rs = append(rs, func(msg tea.WindowSizeMsg) tea.WindowSizeMsg {
			return tea.WindowSizeMsg{
				Width:  (msg.Width / len(row)) - (borderWidth * 2),
				Height: (msg.Height / len(sections)) - (borderWidth * 2),
			}
		})
	}

	return Model{
		sections: sections,
		active:   pos{0, 0},
		resizers: rs,
		Styles:   DefaultStyles(),
		KeyMap:   DefaultKeyMap(),
	}
}

func DefaultStyles() Styles {
	b := lipgloss.RoundedBorder()
	return Styles{
		ActiveStyle:            lipgloss.NewStyle().Border(b),
		InactiveStyle:          lipgloss.NewStyle().Padding(1),
		ActiveBorder:           b,
		HorizontalJoinPosition: lipgloss.Top,
		VerticalJoinPosition:   lipgloss.Top,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		m.switchPane(msg)
	case tea.WindowSizeMsg:
		var cmds []tea.Cmd
		for i := range m.sections {
			for j := range m.sections[i] {
				var c tea.Cmd
				m.sections[i][j], c = m.sections[i][j].Update(
					m.resizers[i](msg),
				)
				cmds = append(cmds, c)
			}
		}
		return m, tea.Batch(cmds...)
	}
	m.sections[m.active.y][m.active.x], cmd = m.sections[m.active.y][m.active.x].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	var rows []string
	for i, row := range m.sections {
		r := []string{}
		for j, cell := range row {
			if m.active.x == j && m.active.y == i {
				r = append(
					r,
					m.Styles.ActiveStyle.Render(
						cell.View(),
					),
				)
			} else {
				r = append(
					r,
					m.Styles.InactiveStyle.Render(
						cell.View(),
					),
				)
			}
		}
		rows = append(
			rows,
			lipgloss.JoinHorizontal(
				m.Styles.HorizontalJoinPosition,
				r...,
			),
		)

	}
	return lipgloss.JoinVertical(
		m.Styles.VerticalJoinPosition,
		rows...,
	)
}

func (m *Model) switchPane(msg tea.KeyMsg) {
	switch {
	case key.Matches(msg, m.KeyMap.Up):
		m.active = m.calcVertical(-1)
	case key.Matches(msg, m.KeyMap.Down):
		m.active = m.calcVertical(1)
	case key.Matches(msg, m.KeyMap.Right):
		m.active = m.calcHorizontal(1)
	case key.Matches(msg, m.KeyMap.Left):
		m.active = m.calcHorizontal(-1)
	}
}

func (m Model) calcVertical(dir int) pos {
	newY := clamp(
		m.active.y+dir,
		0,
		len(m.sections)-1,
	)

	p := pos{
		y: newY,
		x: clamp(m.active.x, 0, len(m.sections[newY])-1),
	}

	return p
}

func (m Model) calcHorizontal(dir int) pos {
	return pos{
		x: clamp(m.active.x+dir, 0, len(m.sections[m.active.y])-1),
		y: m.active.y,
	}
}

func clamp(v int, low int, high int) int {
	return min(high, max(low, v))
}

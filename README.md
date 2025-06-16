# Panes

A [Bubble Tea](https://github.com/charmbracelet/bubbletea) component for 
creating multi-pane applications.

## Default Keybinds

Use ctrl+h, ctrl+j, ctrl+k and ctrl+l to switch between panes, with only the 
active pane being updated.

## Default Styles

By default, the active pane will have a rounded border.

## Describe your pane layout

To describe a pane layout, use a 2D slice of Models like so:

`go
    p := panes.New(
        [][]tea.Model{
        {yourmodel.New(), yourmodel.New()},
        {yourmodel.New(), yourmodel.New()},
        },
    )
`


package components

import (
	"bubbletea-app/app/styles"
)

type HeaderModel struct {
	title string
}

func NewHeaderModel(title string) HeaderModel {
	return HeaderModel{
		title: title,
	}
}

func (m HeaderModel) View(width int) string {
	return styles.HeaderStyle.Width(width - 4).Render(m.title)
}

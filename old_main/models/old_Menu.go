package models

import (
	"fmt"
)

type Menu struct {
	Items [5]MenuOption
}

func (m Menu) PrintItems() {
	for i := 0; i < len(m.Items); i++ {
		fmt.Printf("%d: %s \n", i+1, m.Items[i])
	}
}

func NewMenu() *Menu {
	return &Menu{
		Items: [5]MenuOption{
			Start,
			Add,
			Delete,
			View,
			"EXIT",
		},
	}
}

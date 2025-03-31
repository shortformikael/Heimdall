package container

import "fmt"

type Menu struct {
	Items       *TreeGraph
	Cursor      int
	Current     *TreeNode
	Before      *TreeNode
	Breadcrumbs *LinkedList
}

func NewMenu(pItems *TreeGraph) *Menu {
	var l *LinkedList = &LinkedList{}
	l.Append(pItems.Head)
	return &Menu{
		Items:       pItems,
		Cursor:      0,
		Current:     pItems.Head,
		Breadcrumbs: l,
	}
}

func (m *Menu) PrintCli() {
	fmt.Println("===", m.Current.Value, "===")

	// Breadcrumbs
	//bCrumbs := m.Breadcrumbs.GetArray()
	//fmt.Println(bCrumbs)
	/*
		for _, crumb := range bCrumbs {
			fmt.Print("> ", crumb, " ")
		}*/
	//fmt.Print("\n")

	for i := 0; i < len(m.Current.Children); i++ {
		if i == m.Cursor {
			fmt.Println(" >[*]", m.Current.Children[i])
		} else {
			fmt.Println("  [*]", m.Current.Children[i])
		}
	}
}

func (m *Menu) PrintCliTitle() {
	fmt.Println("===", m.Current.Value, "===")
}

func (m *Menu) Next() {
	if m.Cursor+1 >= len(m.Current.Children) {
		m.Cursor = 0
	} else {
		m.Cursor = m.Cursor + 1
	}
}

func (m *Menu) Previous() {
	if m.Cursor <= 0 {
		m.Cursor = len(m.Current.Children) - 1
	} else {
		m.Cursor = m.Cursor - 1
	}
}

func (m *Menu) Select() string {
	m.Before = m.Current
	m.Current = m.Current.Children[m.Cursor]
	m.Breadcrumbs.Append(m.Current) //????
	m.Cursor = 0
	return string(m.Current.String())
}

func (m *Menu) Back() {
	if m.Before != nil {
		m.Current = m.Before
		m.Before = nil
		m.Cursor = 0
	}
}

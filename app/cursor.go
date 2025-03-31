package app

import (
	"fmt"

	"github.com/shortformikael/GoLearning/container"
)

type cursor struct {
	Menu    *container.TreeNode
	Pointer int
}

func (c *cursor) SetMenu(m *container.TreeNode) {
	c.Menu = m
	fmt.Println("c.Menu = ", c.Menu, c.Menu.Children)
	if c.Menu.Children != nil {
		c.Pointer = 0
	}
}

func (c *cursor) String() string {
	return "Cursor: " + c.Menu.String() + " " + string(c.Pointer)
}

func (c *cursor) Debug() {
	fmt.Println(c.Menu.Children[c.Pointer], c.Menu, c.Pointer, c.Menu.Children)
}

func (c *cursor) Compare(m *container.TreeNode) bool {
	return c.Menu.Children[c.Pointer] == m
}

func (c *cursor) Next() {
	if c.Pointer+1 >= len(c.Menu.Children) {
		c.Pointer = 0
	} else {
		c.Pointer = c.Pointer + 1
	}
}

func (c *cursor) Previous() {
	if c.Pointer-1 < 0 {
		c.Pointer = len(c.Menu.Children) - 1
	} else {
		c.Pointer = c.Pointer - 1
	}
}

func (c *cursor) Select() *container.TreeNode {
	return c.Menu.Children[c.Pointer]
}

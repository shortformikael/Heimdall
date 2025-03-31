package main

import (
	"fmt"

	"github.com/shortformikael/GoLearning/app"
	"github.com/shortformikael/GoLearning/container"
)

var App *app.Engine = &app.Engine{}

func main() {
	fmt.Println("Soft Reset Complete")
	varInit()
	start()
}

func varInit() {
	fmt.Println("Started Program Initialization...")
	menuTree := &container.TreeGraph{Head: &container.TreeNode{}}
	App.List = &container.LinkedList{}

	menuTree.Head.Value = "Main Menu"
	menuTree.Head.AddChild(&container.TreeNode{Value: "View"})
	menuTree.Head.AddChild(&container.TreeNode{Value: "Settings"})
	menuTree.Head.AddChild(&container.TreeNode{Value: "Exit"})

	App.Init(menuTree)
}

func start() {
	App.Start()
}

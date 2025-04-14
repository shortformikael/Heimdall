package main

import (
	"fmt"
	"os"

	"github.com/shortformikael/Heimdall/app"
	"github.com/shortformikael/Heimdall/container"
)

var App *app.Engine = &app.Engine{}

func main() {
	fmt.Println("Soft Reset Complete")
	varInit()
	start()

	os.Exit(0)
}

func varInit() {
	fmt.Println("Started Program Initialization...")
	menuInit()

}

func menuInit() {
	menuTree := &container.TreeGraph{Head: &container.TreeNode{}}
	App.List = &container.LinkedList{}

	menuTree.Head.Value = "Main Menu"
	menuTree.Head.AddChild(&container.TreeNode{Value: "Automation"})
	menuTree.Head.AddChild(&container.TreeNode{Value: "Analysis"})
	menuTree.Head.AddChild(&container.TreeNode{Value: "Capture"})
	menuTree.Head.AddChild(&container.TreeNode{Value: "Settings"})
	menuTree.Head.AddChild(&container.TreeNode{Value: "Exit"})

	App.Init(menuTree)
}

func start() {
	App.Start()
}

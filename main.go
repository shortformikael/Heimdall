package main

import (
	"fmt"
	"os"

	settings "github.com/shortformikael/Heimdall/Settings"
	"github.com/shortformikael/Heimdall/app"
	"github.com/shortformikael/Heimdall/container"
)

var App *app.Engine = &app.Engine{}
var directories = [5]string{
	"./entries",
	"./entries/done",
	"./entries/ongoing",
	"./pcaps",
	"./pcaps/done",
}

func main() {
	fmt.Println("Soft Reset Complete")
	varInit()
	start()

	if settings.Config.Debug {
		cleanTmp()
	}

	os.Exit(0)
}

func cleanTmp() {
	for _, dir := range directories {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			if !file.IsDir() {
				os.Remove(dir + "/" + file.Name())
			}
		}
	}
}

func varInit() {
	fmt.Println("Started Program Initialization...")
	dirInit()
	menuInit()

}

func dirInit() {
	for _, dir := range directories {
		dirCheck(dir)
	}
}

func dirCheck(dir string) {
	_, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Creating Missing Directory:", dir)
		os.Mkdir(dir, 0777)
	}
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

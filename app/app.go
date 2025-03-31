package app

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/shortformikael/GoLearning/container"
)

type Engine struct {
	Menu *container.Menu
	List *container.LinkedList

	Running   bool
	commandCh chan string
	drawCh    chan string
	sigCh     chan os.Signal
	keyCh     chan keyboard.Key
	wg        sync.WaitGroup
}

func (e *Engine) Start() {
	fmt.Println("Starting Engine...")
	e.Running = true
	e.wg.Add(4)
	go e.keyboardListener(0, &e.wg)
	go e.commandListener(1, &e.wg)
	go e.displayListener(2, &e.wg)
	go e.actionListener(3, &e.wg)

	e.wg.Wait()
	fmt.Println("All workers completed")
}

func (e *Engine) Init(tree *container.TreeGraph) {
	e.Running = false
	e.commandCh = make(chan string) //Command Channel,
	e.drawCh = make(chan string)
	e.Menu = container.NewMenu(tree)
	if err := keyboard.Open(); err != nil {
		fmt.Println("Error opening keyboard:", err)
		return
	}
	//defer keyboard.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	e.keyCh = make(chan keyboard.Key)
}

func (e *Engine) Shutdown() {
	fmt.Println("Shutting down engine...")
	os.Exit(1)
}

func (e *Engine) printTreeNode(node *container.TreeNode) {
	fmt.Println("==", node, "==")
	for i := 0; i < len(node.Children); i++ {
		fmt.Println("  [*]", node.Children[i])
	}
}

func (e *Engine) keyboardListener(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Process %d Started\n", id)

	for e.Running {
		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println("Error getting key:", err)
			keyboard.Close()
			fmt.Println("Attempting to re-open keyboard...")

			if err := keyboard.Open(); err != nil {
				fmt.Printf("Failed to re-open keyboard: %v. Exiting... \n", err)
				keyboard.Close()
				e.Shutdown()
				return
			}
			continue
		}
		if key != keyboard.Key(0) {
			e.keyCh <- key
		} else if char != 0 {
			e.keyCh <- keyboard.Key(char)
		}
	}

	fmt.Printf("Process %d Ended\n", id)
}

func (e *Engine) actionListener(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Process %d Started\n", id)

	for e.Running {
		select {
		case <-e.sigCh:
			fmt.Println("\nRecieved interrupt signal. Exiting...")
			e.Shutdown()
			break
		case key, ok := <-e.keyCh:
			if !ok {
				fmt.Println("\nKeyboard channel closed. Exiting...")
				e.Shutdown()
				return
			}

			switch key {
			case keyboard.KeyEsc:
				fmt.Println("\nESC pressed. Exiting...")
				e.Shutdown()

			case keyboard.KeyArrowDown:
				e.commandCh <- "NEXT"
			case keyboard.KeyArrowUp:
				e.commandCh <- "PREVIOUS"
			case keyboard.KeyArrowLeft:
				fmt.Println("[ARROW LEFT]")
			case keyboard.KeyArrowRight:
				fmt.Println("[ARROW RIGHT]")
			case 13:
				e.commandCh <- "SELECT"
			case 8:
				e.commandCh <- "BACK"
			default:
				if key >= 32 && key <= 126 {
					fmt.Printf("Key: %c\n", key)
				} else {
					fmt.Printf("Special Key: %v\n", key)
				}
			}
		}
	}
	fmt.Printf("Process %d Ended\n", id)
}

func (e *Engine) displayListener(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Process %d Started\n", id)

	time.Sleep(1 * time.Second)

	clearConsole()
	e.Menu.PrintCli()

	for e.Running {
		comm := <-e.drawCh
		clearConsole()
		if e.Menu.Current.String() == "Main Menu" {
			e.Menu.PrintCli()
		} else {
			e.Menu.PrintCliTitle()
		}
		fmt.Println("Message:", comm)
	}

	fmt.Printf("Process %d Ended\n", id)
}

func (e *Engine) commandListener(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Process %d Started\n", id)

	for e.Running {
		comm := <-e.commandCh
		switch comm {
		case "NEXT":
			e.Menu.Next()
		case "PREVIOUS":
			e.Menu.Previous()
		case "SELECT":
			sel := e.Menu.Select()
			if sel == "Exit" {
				e.Shutdown()
			} else {
				e.drawCh <- sel
			}
			continue
		case "BACK":
			e.Menu.Back()
		}
		e.drawCh <- ""
	}
}

func clearConsole() {
	fmt.Println(("\033[H\033[2J"))
}

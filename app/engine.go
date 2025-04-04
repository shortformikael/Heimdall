package app

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/shortformikael/Heimdall/analyzer"
	"github.com/shortformikael/Heimdall/capturer"
	"github.com/shortformikael/Heimdall/container"
	"github.com/shortformikael/Heimdall/sender"
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

	capturer *capturer.CaptureManager
	analyzer *analyzer.AnalyzerManager
	sender   *sender.Sender
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
	fmt.Println("All Listeners completed")
}

func (e *Engine) Init(tree *container.TreeGraph) {
	e.Running = false

	e.commandCh = make(chan string) //Command Channel,
	e.drawCh = make(chan string)
	e.sigCh = make(chan os.Signal, 1)
	signal.Notify(e.sigCh, os.Interrupt, syscall.SIGTERM)
	e.keyCh = make(chan keyboard.Key)

	e.Menu = container.NewMenu(tree)
	e.capturer = &capturer.CaptureManager{}
	e.analyzer = &analyzer.AnalyzerManager{}
	e.sender = sender.NewSender()
	e.capturer.Init(e.drawCh)
	e.analyzer.Init(e.drawCh)

	if err := keyboard.Open(); err != nil {
		fmt.Println("Error opening keyboard:", err)
		return
	}
	//defer keyboard.Close()

}

func (e *Engine) Shutdown() {
	fmt.Println("Shutting down engine...")
	close(e.sigCh)
	e.Running = false
	keyboard.Close()
}

func (e *Engine) keyboardListener(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer fmt.Printf("Process %d Ended\n", id)
	fmt.Printf("Process %d Started\n", id)
	for e.Running {
		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println("Error getting key:", err)
			keyboard.Close()
			fmt.Println("Attempting to re-open keyboard...")
			if err := keyboard.Open(); err != nil {
				fmt.Printf("Failed to re-open keyboard: %v. Exiting... \n", err)
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
}

func (e *Engine) actionListener(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer fmt.Printf("Process %d Ended\n", id)
	fmt.Printf("Process %d Started\n", id)
	for {
		select {
		case <-e.sigCh:
			fmt.Println("actionListerner Shutdown")
			return
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
			case 13:
				e.commandCh <- "SELECT"
			case 8, 127:
				e.commandCh <- "BACK"
			case keyboard.KeyArrowLeft:
				fmt.Println("[ARROW LEFT]")
			case keyboard.KeyArrowRight:
				fmt.Println("[ARROW RIGHT]")
			default:
				if key >= 32 && key <= 126 {
					fmt.Printf("Key: %c\n", key)
				} else {
					fmt.Printf("Special Key: %v\n", key)
				}
			}
		}
	}
}

func (e *Engine) displayListener(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer fmt.Printf("Process %d Ended\n", id)
	fmt.Printf("Process %d Started\n", id)

	time.Sleep(1 * time.Second)

	clearConsole()
	fmt.Println("")
	e.Menu.PrintCli()

	for {
		select {
		case <-e.sigCh:
			fmt.Println("DisplayListener Shutdown")
			return
		case comm := <-e.drawCh:
			clearConsole()
			if comm != "" {
				fmt.Println("Event:", comm)
			} else {
				fmt.Println("")
			}

			switch e.Menu.Current.String() {
			case "Main Menu":
				e.Menu.PrintCli()
			case "Automation":
				fmt.Println("=== Automation ===")
				e.capturer.PrintAutomation()
				e.analyzer.PrintAutomation()
				e.sender.PrintAutomation()
			case "Capture":
				e.Menu.PrintCliTitle()
				e.capturer.PrintCli()
			case "Analysis":
				e.analyzer.PrintCli()
			default:
				e.Menu.PrintCliTitle()
			}
		}
	}
}

func (e *Engine) commandListener(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer fmt.Printf("Process %d Ended \n", id)
	fmt.Printf("Process %d Started\n", id)

	for {

		select {
		case <-e.sigCh:
			fmt.Println("commandListener Shutdown")
			return
		case comm := <-e.commandCh:
			switch comm {
			case "NEXT":
				e.Menu.Next()
			case "PREVIOUS":
				e.Menu.Previous()
			case "SELECT":
				switch e.Menu.Current.String() {
				case "Main Menu":
					sel := e.Menu.Select()
					if sel == "Exit" {
						go e.Shutdown()
						continue
					} else {
						e.drawCh <- "Selected " + sel
						continue
					}
				case "Automation":
					if !e.analyzer.Running {
						e.analyzer.StartAutomation()
					} else {
						e.analyzer.EndAutomation()
					}
					if !e.capturer.Running {
						e.capturer.StartAutomation()
					} else {
						e.capturer.EndAutomation()
					}
					if !e.sender.Running {
						e.sender.StartAutomation()
					} else {
						e.sender.EndAutomation()
					}
				case "Capture":
					if e.capturer.Running {
						e.capturer.EndCapture()
						e.drawCh <- "End Capture"
					} else {
						err := e.capturer.StartCapture()
						if err != nil {
							e.drawCh <- err.Error()
							continue
						}
						e.drawCh <- "Start Capture"
					}
					continue
				case "Analysis":
					e.analyzer.Start()
					continue
				}
			case "BACK":
				e.Menu.Back()
			}
			e.drawCh <- ""
		}
	}
}

func clearConsole() {
	// fmt.Println(("\033[H\033[2J"))

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

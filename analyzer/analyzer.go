package analyzer

import (
	"fmt"
	"os"
	"sync"
)

type AnalyzerManager struct {
	Running         bool
	dirPath         string
	conversationMap map[string]*Conversation
	drawCh          *chan string
	availableFiles  []string
	wg              *sync.WaitGroup
}

func (a *AnalyzerManager) Init(drawCh *chan string) {
	a.conversationMap = make(map[string]*Conversation)
	a.drawCh = drawCh
	a.availableFiles = a.getAvailableFiles()
	a.wg = &sync.WaitGroup{}
	a.dirPath = "./pcaps"
}

func (a *AnalyzerManager) StartAutomation() {
	a.Running = true
}

func (a *AnalyzerManager) EndAutomation() {
	a.Running = false
}

func (a *AnalyzerManager) Start() {
	a.availableFiles = a.getAvailableFiles()
	//Iterate through available files
	for _, file := range a.availableFiles {
		//Create worker for each file
		a.wg.Add(1)
		worker := NewWorker(a.wg, file)
		go worker.Start()
		//Have each worker serialize a jsonfile
	}
	a.wg.Wait()
	fmt.Println("Done!")
}

func (a *AnalyzerManager) getAvailableFiles() []string {
	r := []string{}
	entries, err := os.ReadDir(a.dirPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			r = append(r, entry.Name())
		}
	}

	return r
}

func (a *AnalyzerManager) PrintAutomation() {
	if a.Running {
		fmt.Println(" -> Analyze Running...")
	} else {
		fmt.Println(" -> Waiting to Analyze...")
	}
	fmt.Println(" -> ")
	fmt.Println(" -> ")
	fmt.Println(" -> ")
}

func (a *AnalyzerManager) PrintCli() {
	fmt.Println("=== Analysis ===")
	fmt.Println(" -> Youre in the Analyzer menu")
	a.PrintAvailableFiles()
}

func (a *AnalyzerManager) PrintAvailableFiles() {
	for _, files := range a.getAvailableFiles() {
		fmt.Println("-", files)
	}
}

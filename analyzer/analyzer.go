package analyzer

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

type AnalyzerManager struct {
	Running         bool
	targetPath      string
	ongoingPath     string
	donePath        string
	conversationMap map[string]*Conversation
	availableFiles  []string
	wg              *sync.WaitGroup

	drawCh    chan string
	SigCh     chan string
	AnalyzeCh chan string

	track1 string
	track2 string
	track3 string
}

func (a *AnalyzerManager) Init(drawCh chan string) {
	a.conversationMap = make(map[string]*Conversation)

	a.drawCh = drawCh
	a.AnalyzeCh = make(chan string)
	a.SigCh = make(chan string)

	a.availableFiles = a.getAvailableFiles()
	a.wg = &sync.WaitGroup{}
	a.targetPath = "./pcaps/done"
	a.ongoingPath = "./entries/ongoing"
	a.donePath = "./entries"
	a.track1 = ""
	a.track2 = " - Running: "
	a.track3 = ""
}

func (a *AnalyzerManager) StartAutomation() {
	a.Running = true
	a.updateTrackers(
		"update",
		"",
		"",
	)

}

func (a *AnalyzerManager) actionListener() {
	for {
		select {
		case <-a.SigCh:
			return
		case <-a.AnalyzeCh:
			a.Start()
		}
	}
}

func (a *AnalyzerManager) EndAutomation() {
	a.updateTrackers(
		"update",
		"",
		"",
	)
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
	entries, err := os.ReadDir(a.targetPath)
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
		fmt.Println(" =====> ANALYZE RUNNING...")
	} else {
		fmt.Println(" =====> not analyzing")
	}
	fmt.Println(a.track1)
	fmt.Println(a.track2)
	fmt.Println(a.track3)
}

func (a *AnalyzerManager) updateTrackers(t1 string, t2 string, t3 string) {
	if t1 != "" {
		a.track1 = " - Total Available: " + strconv.FormatInt(int64(len(a.getAvailableFiles())), 10)
	}
	if t2 != "" {
		a.track2 = t2
	}
	if t3 != "" {
		a.track3 = t3
	}
	a.drawCh <- ""
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

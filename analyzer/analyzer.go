package analyzer

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type AnalyzerManager struct {
	Running         bool
	targetPath      string
	ongoingPath     string
	donePath        string
	conversationMap map[string]*Conversation
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

	a.wg = &sync.WaitGroup{}
	a.targetPath = "./pcaps/done"
	a.ongoingPath = "./entries/ongoing"
	a.donePath = "./entries/done"
	a.track1 = ""
	a.track2 = ""
	a.track3 = ""

	//a.updateTrackers("update", " - Running: 0", "")
}

func (a *AnalyzerManager) StartAutomation() {
	a.Running = true
	a.updateTrackers(
		"update",
		" - Running:",
		"",
	)
	go a.actionListener()
	go a.StartAuto()
}

func (a *AnalyzerManager) actionListener() {
	for {
		select {
		case <-a.SigCh:
			return
		case <-a.AnalyzeCh:
			go a.StartAuto()
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

func (a *AnalyzerManager) StartAuto() {
	availbleFiles := a.getAvailableFiles()
	count := 0

	for _, file := range availbleFiles {

		//Move the file to ongoing
		split := strings.Split(file, "/")
		dstPath := a.ongoingPath + "/" + split[len(split)-1]
		err := os.Rename(file, dstPath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		count++
		a.wg.Add(1)
		worker := NewWorker(a.wg, dstPath, a.donePath)
		go worker.Start()
	}
	a.updateTrackers("update", "Running: "+strconv.FormatInt(int64(count), 10), "")

	a.wg.Wait()
	a.updateTrackers("update", "", " - Done with some workers")

	time.Sleep(1 * time.Second)
	if a.Running {
		a.AnalyzeCh <- ""
	}
}

func (a *AnalyzerManager) Start() {
	availableFiles := a.getAvailableFiles()
	//Iterate through available files an create worker for each
	for _, file := range availableFiles {

		split := strings.Split(file, "/")
		dstPath := a.ongoingPath + "/" + split[len(split)-1]
		err := os.Rename(file, dstPath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//Create worker for each file
		a.wg.Add(1)
		worker := NewWorker(a.wg, dstPath, a.donePath)
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
			r = append(r, a.targetPath+"/"+entry.Name())
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

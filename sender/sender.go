package Sender

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Sender struct {
	Running  bool
	pcapPath string
	jsonPath string

	sigCh    chan string
	SenderCh chan string
	count    int
}

func NewSender() *Sender {
	return &Sender{
		Running:  false,
		pcapPath: "./entries/done",
		jsonPath: "./entries",
		SenderCh: make(chan string),
		count:    0,
	}
}

func (r *Sender) StartAutomation() {
	r.Running = true

	go r.actionListener()
	go r.StartReading()
}

func (r *Sender) actionListener() {
	r.sigCh = make(chan string)
	for {
		select {
		case <-r.sigCh:
			return
		case <-r.SenderCh:
			go r.StartReading()
		}
	}
}

func (r *Sender) StartReading() {
	pcapFiles := r.getAvailableFiles(r.pcapPath)
	jsonFiles := r.getAvailableFiles(r.jsonPath)
	// Remove processed pcaps
	for _, file := range pcapFiles {
		os.Remove(file)
	}
	// Send and Remove json files
	for _, file := range jsonFiles {
		r.count++
		os.Remove(file)
	}

	time.Sleep(1 * time.Second)
	if r.Running {
		r.SenderCh <- ""
	}
}

func (r *Sender) getAvailableFiles(path string) []string {
	ret := []string{}
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			ret = append(ret, (path + "/" + entry.Name()))
		}
	}

	return ret
}

func (r *Sender) EndAutomation() {
	r.Running = false
	close(r.sigCh)
}

func (r *Sender) PrintAutomation() {
	if r.Running {
		fmt.Println(" =====> READING...")
	} else {
		fmt.Println(" =====> not reading")
	}

	fmt.Println(" - Total files sent:", strconv.FormatInt(int64(r.count), 10))
}

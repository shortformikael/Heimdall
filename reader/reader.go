package reader

import (
	"fmt"
	"os"
	"time"
)

type Reader struct {
	Running  bool
	pcapPath string
	jsonPath string

	sigCh    chan string
	ReaderCh chan string
}

func NewReader() *Reader {
	return &Reader{
		Running:  false,
		pcapPath: "./entries/done",
		jsonPath: "./entries",
		ReaderCh: make(chan string),
	}
}

func (r *Reader) StartAutomation() {
	r.Running = true

	go r.actionListener()

	go r.StartReading()
}

func (r *Reader) actionListener() {
	r.sigCh = make(chan string)
	for {
		select {
		case <-r.sigCh:
			return
		case <-r.ReaderCh:
			go r.StartReading()
		}
	}
}

func (r *Reader) StartReading() {
	pcapFiles := r.getAvailableFiles(r.pcapPath)
	jsonFiles := r.getAvailableFiles(r.jsonPath)
	for _, file := range pcapFiles {
		fmt.Println(file)
	}

	for _, file := range jsonFiles {
		fmt.Println(file)
	}

	time.Sleep(1 * time.Second)
	if r.Running {
		r.ReaderCh <- ""
	}
}

func (r *Reader) getAvailableFiles(path string) []string {
	ret := []string{}
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			ret = append(ret, entry.Name())
		}
	}

	return ret
}

func (r *Reader) EndAutomation() {
	r.Running = false
	close(r.sigCh)
}

func (r *Reader) PrintAutomation() {
	if r.Running {
		fmt.Println(" =====> READING...")
	} else {
		fmt.Println(" =====> not reading")
	}
}

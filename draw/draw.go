package draw

import "fmt"

var Channel chan string = make(chan string)

func Start() {
	for {
		draw()
		msg := <-Channel
		fmt.Print(msg)
	}
}

func draw() {

}

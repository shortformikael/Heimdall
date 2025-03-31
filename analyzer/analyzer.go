package analyzer

import "fmt"

type Analyzer struct {
}

func (a *Analyzer) PrintCli() {
	fmt.Println(" -> Youre in the Analyzer menu")
}

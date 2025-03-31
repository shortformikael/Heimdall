package old_main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/shortformikael/GoLearning/old_main/models"
)

var reader *bufio.Reader = bufio.NewReader(os.Stdin)
var menu_Start *models.Menu = models.NewMenu()
var menuStatus models.MenuOption = models.Start

func old_main() {

	fmt.Println("=== START ===")
	Run()
	fmt.Println("=== END ===")
}

func Run() {
	for {
		fmt.Printf("Current Status: %s \n", menuStatus)

		if string(menuStatus) == "EXIT" {
			break
		}
		printMenu()
		SelectChoice(Choose())
	}
}

func printMenu() {
	switch string(menuStatus) {
	case string(models.Start):
		menu_Start.PrintItems()
	case string(models.Add):
		//ADD
	case string(models.Delete):
		//DELETE
	case string(models.View):
		//VIEW
	case string("EXIT"):
		//EXIT
	default:
		//default
	}

}

func Choose() interface{} {
	input := ReadStdInput()
	num, err := strconv.Atoi(input)
	if err != nil {
		return strings.ToUpper(input)
	} else if 1 <= num && num <= len(menu_Start.Items) {
		return num
	}
	return "Out of bounds ERROR" // Return invalid input
}

func SelectChoice(pChoice interface{}) {
	switch string(menuStatus) {
	case string(models.Start):
		pType := reflect.TypeOf(pChoice).Kind()
		if pType == reflect.Int || pType == reflect.String {
			switch pChoice := pChoice.(type) {
			case int:
				menuStatus = menu_Start.Items[pChoice-1]
			case string:
				for item := range menu_Start.Items {
					if string(menu_Start.Items[item]) == pChoice {
						menuStatus = menu_Start.Items[item]
					}
				}

			}
		}
	case string(models.Add):
		//ADD
	case string(models.Delete):
		//DELETE
	case string(models.View):
		//VIEW
	case string("EXIT"):
		//EXIT
	default:
		//default
	}

}

func ReadStdInput() string {
	fmt.Print("Choose option: ")
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

package display

import "fmt"

func DisplayError(errorStr string) {
	fmt.Println("## ОШИБКА!!! ##")
	fmt.Printf("%s\n", errorStr)
}

func DisplaySuccess(successStr string) {
	fmt.Printf("%s\n", successStr)
}

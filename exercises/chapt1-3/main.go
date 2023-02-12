package main

import "fmt"

func main() {
	array := [5]string{"I", "am", "stupid", "and", "weak"}

	for i := 0; i < len(array); i++ {
		switch array[i] {
		case "stupid":
			array[i] = "smart"
		case "weak":
			array[i] = "strong"
		}
	}

	fmt.Println(array)
}

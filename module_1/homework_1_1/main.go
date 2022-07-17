package main

import "fmt"

func main() {
	str := []string{"I", "am", "stupid", "and", "weak"}

	replaceMap := map[int]string{
		2: "smart",
		4: "strong",
	}

	for k := range str {
		replaceWord, ok := replaceMap[k]
		if ok && replaceWord != "" {
			str[k] = replaceWord
		}
	}

	fmt.Printf("str: %+v\n", str)
}

package main

import "fmt"

func main() {
	process, err := GetAllProcess()
	if err != nil {
		panic(err)
	}
	for _, v := range process {
		fmt.Println(v)
	}
}

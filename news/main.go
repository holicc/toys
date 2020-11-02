package main

import "fmt"

func main() {
	news, err := GetNews()
	if err != nil {
		panic(err)
	}
	fmt.Println(news)
}

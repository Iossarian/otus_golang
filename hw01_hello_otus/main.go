package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	greetings := stringutil.Reverse("Hello, OTUS!")
	fmt.Println(greetings)
}

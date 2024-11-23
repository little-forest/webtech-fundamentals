package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.Handle("/",
		http.FileServer(http.Dir("static")))

	port := getPortNumber()
	fmt.Printf("listening port : %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

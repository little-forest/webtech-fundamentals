package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) { // <3>
	fmt.Fprint(w, "Hello, Web application!") // <4>
}

func main() {
	http.HandleFunc("/", hello)              // <2>
	err := http.ListenAndServe(":8080", nil) // <1>
	if err != nil {
		log.Fatal("failed to start : ", err)
	}
}

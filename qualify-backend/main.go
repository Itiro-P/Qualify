package main

import (
	"fmt"
	"net/http"

	"golang.org/x/example/hello/reverse"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Dockerized Go Application!")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is running on port 8001")
	fmt.Println(reverse.String("Hello"))
	http.ListenAndServe(":8001", nil)
}

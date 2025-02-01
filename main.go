package main

import (
	"net/http"

	"github.com/BMokarzel/multithread.git/service"
)

func main() {

	http.HandleFunc("/", service.GetAddressHandler)
	http.ListenAndServe(":8080", nil)

}

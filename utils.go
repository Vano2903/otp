package main

import (
	"fmt"
	"net/http"
)

//printInternalErr set the status code to 500 of the http response
func PrintInternalErr(w http.ResponseWriter, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s"}`, err)))
}

//printErr will return 400 error code to the client
func PrintErr(w http.ResponseWriter, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf(`{"code": 400, "msg": "%s"}`, err)))
}

//return the unavailable service response (adding a unaviable page)
func UnavailablePage(w http.ResponseWriter) {
	//TODO add an unavailable page
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte("{\"msg\": \"page unavailable at the moment\"}"))
}

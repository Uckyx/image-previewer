package main

import (
	"fmt"
	"net/http"
)

//Описать все роуты, возможно нужен отдельный клиент который будет содержать в себе все.

func main() {
	http.HandleFunc("/test", handler) // each request calls handler

	http.ListenAndServe(":8090", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

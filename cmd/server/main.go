package main

import (
	"cyber-security-minigames/views/pages"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	component := pages.Landing()
	http.Handle("/", templ.Handler(component))
	fmt.Println("listening on port 3000")
	http.ListenAndServe(":3000", nil)
}

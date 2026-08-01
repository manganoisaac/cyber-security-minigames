package main

import (
	"fmt"
	"net/http"
	"os"

	"cyber-security-minigames/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Println("listening on port " + port)
	if err := http.ListenAndServe(":"+port, handlers.Routes()); err != nil {
		fmt.Println("server error:", err)
	}
}

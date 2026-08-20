package main

import (
	"fmt"
	"net/http"
	"os"

	"cyber-security-minigames/internal/handlers"
	"cyber-security-minigames/internal/handlers/email_client"
	"cyber-security-minigames/internal/handlers/encryption"
	"cyber-security-minigames/internal/handlers/password"
	"cyber-security-minigames/internal/session"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	store := session.Store()
	games := []handlers.Game{
		email_client.NewGame(email_client.NewHandler(email_client.EmailService())),
		password.NewGame(password.NewHandler(password.PasswordService())),
		encryption.NewGame(encryption.NewHandler()),
	}
	handler := session.Middleware(store)(handlers.Routes(games...))

	fmt.Println("listening on port " + port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		fmt.Println("server error:", err)
	}
}

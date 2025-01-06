package router

import (
	"net/http"

	"go-api/handlers"
	"go-api/middleware"
)

// SetupRouter thiết lập tất cả các route
func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Định nghĩa các route
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.Handle("/user", middleware.Authenticate(http.HandlerFunc(handlers.GetUserHandler)))

	return mux
}

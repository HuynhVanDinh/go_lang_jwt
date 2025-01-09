// package router

// import (
// 	"net/http"

// 	"go-api/handlers"
// 	"go-api/middleware"
// )

// // SetupRouter thiết lập tất cả các route
// func SetupRouter() *http.ServeMux {
// 	mux := http.NewServeMux()

// 	// Định nghĩa các route
// 	mux.HandleFunc("/register", handlers.RegisterHandler)
// 	mux.HandleFunc("/login", handlers.LoginHandler)
// 	mux.Handle("/user", middleware.Authenticate(http.HandlerFunc(handlers.GetUserHandler)))

//		return mux
//	}
package routers

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"go-api/handler"
	"go-api/middleware"
)

// SetupRouter thiết lập tất cả các route
func SetupRouter() http.Handler {
	r := mux.NewRouter()

	// CORS middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}), // Cho phép frontend truy cập
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Public routes
	r.HandleFunc("/register", handler.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handler.LoginHandler).Methods("POST")

	// Protected routes (cần JWT)
	protectedRoutes := r.PathPrefix("/").Subrouter()
	protectedRoutes.Use(middleware.Authenticate) // Middleware áp dụng cho tất cả route bên dưới
	protectedRoutes.HandleFunc("/user", handler.GetUserHandler).Methods("GET")

	return corsHandler(r)
}

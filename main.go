package main

import (
	"log"
	"net/http"

	"go-api/db"
	router "go-api/routers"
)

func main() {
	// Kết nối cơ sở dữ liệu
	err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Không thể kết nối CSDL: %v", err)
	}

	// Định nghĩa các route
	mux := router.SetupRouter()

	// Chạy server
	log.Println("Server đang chạy tại http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

package handler

import (
	"encoding/json"

	"net/http"

	"go-api/db"
	"go-api/models"
	"go-api/utils"

)



// GetUserHandler trả về thông tin người dùng
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Không thể xác thực người dùng", http.StatusUnauthorized)
		return
	}

	// Truy vấn người dùng
	var user models.User
	query := "SELECT id, username, email FROM users WHERE id = ?"
	err := db.DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		http.Error(w, "Người dùng không tồn tại", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"go-api/db"
	"go-api/models"
	"go-api/utils"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Dữ liệu không hợp lệ", http.StatusBadRequest)
		return
	}

	// Hash mật khẩu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Không thể mã hóa mật khẩu", http.StatusInternalServerError)
		return
	}

	// Lưu người dùng vào cơ sở dữ liệu
	_, err = db.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Không thể lưu người dùng", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Đăng ký thành công"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Dữ liệu không hợp lệ", http.StatusBadRequest)
		return
	}

	var user models.User
	err = db.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", req.Username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		http.Error(w, "Sai tên đăng nhập hoặc mật khẩu", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Lỗi hệ thống", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Sai tên đăng nhập hoặc mật khẩu", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Không thể tạo token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

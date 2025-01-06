package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"go-api/db"
	"go-api/models"
	"go-api/utils"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler xử lý đăng ký người dùng
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Chỉ hỗ trợ phương thức POST", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Dữ liệu không hợp lệ", http.StatusBadRequest)
		return
	}

	// Kiểm tra dữ liệu đầu vào
	if user.Username == "" || user.Password == "" || user.Email == "" {
		http.Error(w, "Tên người dùng, mật khẩu và email không được để trống", http.StatusBadRequest)
		return
	}

	// Hash mật khẩu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Không thể mã hóa mật khẩu", http.StatusInternalServerError)
		return
	}

	// Lưu vào cơ sở dữ liệu
	_, err = db.DB.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?)", user.Username, string(hashedPassword), user.Email)
	if err != nil {
		if sqlErr, ok := err.(*mysql.MySQLError); ok && sqlErr.Number == 1062 {
			http.Error(w, "Tên người dùng hoặc email đã tồn tại", http.StatusConflict)
		} else {
			http.Error(w, "Không thể lưu người dùng", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Đăng ký thành công")
}

// LoginHandler xử lý đăng nhập
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Dữ liệu không hợp lệ", http.StatusBadRequest)
		return
	}

	// Tìm người dùng trong cơ sở dữ liệu
	var user models.User
	query := "SELECT id, username, password FROM users WHERE username = ?"
	err = db.DB.QueryRow(query, req.Username).Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		http.Error(w, "Sai tên đăng nhập hoặc mật khẩu", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Lỗi hệ thống", http.StatusInternalServerError)
		return
	}

	// Kiểm tra mật khẩu
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Sai tên đăng nhập hoặc mật khẩu", http.StatusUnauthorized)
		return
	}

	// Tạo JWT
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Không thể tạo token", http.StatusInternalServerError)
		return
	}

	// Trả về token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// GetUserHandler trả về thông tin người dùng
// func GetUserHandler(w http.ResponseWriter, r *http.Request) {
// 	userID, ok := utils.GetUserID(r.Context())
// 	if !ok {
// 		http.Error(w, "Không thể xác thực người dùng", http.StatusUnauthorized)
// 		return
// 	}

// 	// Truy vấn người dùng
// 	var user models.User
// 	query := "SELECT id, username, email FROM users WHERE id = ?"
// 	err := db.DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email)
// 	if err != nil {
// 		http.Error(w, "Người dùng không tồn tại", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(user)
// }

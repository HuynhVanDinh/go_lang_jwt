package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

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

// func LoginHandler(w http.ResponseWriter, r *http.Request) {
// 	var req models.LoginRequest
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		http.Error(w, "Dữ liệu không hợp lệ", http.StatusBadRequest)
// 		return
// 	}

// 	var user models.User
// 	err = db.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", req.Username).
// 		Scan(&user.ID, &user.Username, &user.Password)
// 	if err == sql.ErrNoRows {
// 		http.Error(w, "Sai tên đăng nhập hoặc mật khẩu", http.StatusUnauthorized)
// 		return
// 	} else if err != nil {
// 		http.Error(w, "Lỗi hệ thống", http.StatusInternalServerError)
// 		return
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
// 	if err != nil {
// 		http.Error(w, "Sai tên đăng nhập hoặc mật khẩu", http.StatusUnauthorized)
// 		return
// 	}

// 	token, err := utils.GenerateToken(user.ID)
// 	if err != nil {
// 		http.Error(w, "Không thể tạo token", http.StatusInternalServerError)
// 		return
// 	}

//		json.NewEncoder(w).Encode(map[string]string{"token": token})
//	}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Dữ liệu không hợp lệ"}`, http.StatusBadRequest)
		return
	}

	var user models.User
	err := db.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", req.Username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error": "Sai tên đăng nhập hoặc mật khẩu"}`, http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("Lỗi truy vấn database: %v", err)
		http.Error(w, `{"error": "Lỗi hệ thống"}`, http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, `{"error": "Sai tên đăng nhập hoặc mật khẩu"}`, http.StatusUnauthorized)
		return
	}
	// Xóa tất cả các bản ghi cũ trong login_history
	// _, err = db.DB.Exec("DELETE FROM login_history WHERE user_id = ? AND id NOT IN (SELECT id FROM login_history WHERE user_id = ? ORDER BY login_time DESC LIMIT 1)", user.ID)
	// if err != nil {
	// 	http.Error(w, "Lỗi khi xóa lịch sử đăng nhập cũ", http.StatusInternalServerError)
	// 	return
	// }
	// Ghi lịch sử đăng nhập
	ipAddress := strings.Split(r.RemoteAddr, ":")[0] // Lấy IP từ request
	deviceInfo := r.Header.Get("User-Agent")         // Thông tin thiết bị từ User-Agent

	// Kiểm tra đăng nhập từ thiết bị mới
	var lastDevice string
	err = db.DB.QueryRow("SELECT device_info FROM login_history WHERE user_id = ? ORDER BY login_time DESC LIMIT 1", user.ID).
		Scan(&lastDevice)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Lỗi kiểm tra thiết bị trước đó: %v", err)
		http.Error(w, `{"error": "Lỗi hệ thống"}`, http.StatusInternalServerError)
		return
	}

	isNewDevice := lastDevice != "" && lastDevice != deviceInfo

	_, err = db.DB.Exec("INSERT INTO login_history (user_id, device_info, ip_address) VALUES (?, ?, ?)",
		user.ID, deviceInfo, ipAddress)
	if err != nil {
		log.Printf("Lỗi ghi lịch sử đăng nhập: %v", err)
		http.Error(w, `{"error": "Lỗi hệ thống"}`, http.StatusInternalServerError)
		return
	}

	// Tạo JWT
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		log.Printf("Lỗi tạo token: %v", err)
		http.Error(w, `{"error": "Không thể tạo token"}`, http.StatusInternalServerError)
		return
	}

	// Cảnh báo nếu đăng nhập từ thiết bị mới
	if isNewDevice {
		log.Printf("Cảnh báo: Người dùng %s đăng nhập từ thiết bị mới!", user.Username)
		// Gửi email hoặc thông báo tại đây (nếu cần)
	}

	// Trả về token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Đăng nhập thành công",
		"token":   token,
	})
}
func GetLoginHistory(w http.ResponseWriter, r *http.Request) {
	// Lấy user ID từ query params
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, `{"error": "Thiếu user_id"}`, http.StatusBadRequest)
		return
	}

	// Truy vấn lịch sử đăng nhập
	rows, err := db.DB.Query(`
		SELECT login_time, device_info, ip_address 
		FROM login_history 
		WHERE user_id = ? 
		ORDER BY login_time DESC
	`, userID)
	if err != nil {
		log.Printf("Lỗi truy vấn lịch sử đăng nhập: %v", err)
		http.Error(w, `{"error": "Lỗi hệ thống"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Duyệt qua kết quả và đưa vào danh sách
	var history []map[string]interface{}
	for rows.Next() {
		var loginTime time.Time
		var deviceInfo, ipAddress string

		if err := rows.Scan(&loginTime, &deviceInfo, &ipAddress); err != nil {
			log.Printf("Lỗi scan dữ liệu: %v", err)
			http.Error(w, `{"error": "Lỗi hệ thống"}`, http.StatusInternalServerError)
			return
		}

		history = append(history, map[string]interface{}{
			"login_time":  loginTime.Format(time.RFC3339), // Format ngày giờ chuẩn ISO
			"device_info": deviceInfo,
			"ip_address":  ipAddress,
		})
	}

	// Kiểm tra nếu không có lịch sử
	if len(history) == 0 {
		http.Error(w, `{"error": "Không có lịch sử đăng nhập"}`, http.StatusNotFound)
		return
	}

	// Trả về JSON kết quả
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Lịch sử đăng nhập",
		"history": history,
	})
}
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Không xác định được user", http.StatusUnauthorized)
		return
	}

	// Lấy token từ tiêu đề Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, `{"error": "Token không hợp lệ hoặc không tồn tại"}`, http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Kiểm tra token có hợp lệ không (nếu lưu token trong cơ sở dữ liệu, kiểm tra ở đây)
	_, err := utils.VerifyToken(token)
	if err != nil {
		http.Error(w, `{"error": "Token không hợp lệ"}`, http.StatusUnauthorized)
		return
	}

	// Lưu lịch sử đăng xuất (nếu cần thiết)
	ipAddress := strings.Split(r.RemoteAddr, ":")[0]
	deviceInfo := r.Header.Get("User-Agent")

	_, err = db.DB.Exec("INSERT INTO logout_history (user_id,token, ip_address, device_info, logout_time) VALUES (?, ?, ?, ?, ?)",
		userID, token, ipAddress, deviceInfo, time.Now())
	if err != nil {
		log.Printf("Lỗi ghi lịch sử đăng xuất: %v", err)
		http.Error(w, `{"error": "Lỗi hệ thống"}`, http.StatusInternalServerError)
		return
	}

	// Phản hồi thành công
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Đăng xuất thành công",
	})
}

// package db

// import (
// 	"database/sql"
// 	"log"

// 	_ "github.com/go-sql-driver/mysql"
// )

// var DB *sql.DB

// // InitDB khởi tạo kết nối MySQL
// func InitDB() {
// 	var err error
// 	dsn := "root:123456@tcp(127.0.0.1:3307)/golang_demo"
// 	DB, err = sql.Open("mysql", dsn)
// 	if err != nil {
// 		log.Fatalf("Không thể mở kết nối MySQL: %v", err)
// 	}

// 	if err = DB.Ping(); err != nil {
// 		log.Fatalf("Không thể kết nối MySQL: %v", err)
// 	}

//		log.Println("Kết nối MySQL thành công!")
//	}
package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() error {
	var err error
	DB, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3307)/golang_demo?parseTime=true")
	if err != nil {
		return err
	}

	// Kiểm tra kết nối
	err = DB.Ping()
	if err != nil {
		return err
	}

	log.Println("Kết nối cơ sở dữ liệu thành công")
	return nil
}
func CheckLogoutHistory(token string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM logout_history WHERE token = ?)"
	err := DB.QueryRow(query, token).Scan(&exists)
	if err != nil {
		log.Printf("Lỗi khi truy vấn database: %v", err)
		return false, err
	}
	return exists, nil
}


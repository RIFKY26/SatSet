package domain

// Menggunakan nama AuthUser agar tidak tertukar dengan User Service
type User struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

// FUNGSI PENTING: Memberi tahu GORM agar membuat tabel bernama "auth_users"
func (User) TableName() string {
	return "auth_users"
}
package domain

// Menggunakan nama AuthUser agar tidak tertukar dengan User Service
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

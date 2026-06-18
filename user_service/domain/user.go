package domain

import "time"

type User struct {
	ID        string    `json:"id" gorm:"primaryKey"` // Tambahkan gorm:"primaryKey"
	Username  string    `json:"username"`
	Email     string    `json:"email" gorm:"unique"`  // Tambahkan gorm:"unique" agar email tidak dobel
	Password  string    `json:"-"` 
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
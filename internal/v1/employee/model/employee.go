package model

import "time"

type Employee struct {
	ID        int       `json:"id"`
	FullName  string    `json:"fullname"`
	Salary    float64   `json:"salary"`
	Code      string    `json:"code"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
}

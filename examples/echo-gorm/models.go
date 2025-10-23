package main

import (
	"time"

	"github.com/lucasvillarinho/restql"
)

// User represents a user in the system
type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"not null;uniqueIndex" json:"email"`
	Status    string    `gorm:"not null" json:"status"`
	Age       int       `gorm:"not null" json:"age"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName overrides the table name used by User
func (User) TableName() string {
	return "users"
}

// UsersSchema is the RestQL schema for users table
var UsersSchema = restql.NewSchema("users").
	AllowFields("id", "name", "email", "status", "age", "created_at")

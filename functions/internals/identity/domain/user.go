package domain

import (
"errors"
"time"
)

// Role defines the user's permission level
type Role string

const (
RoleAdmin  Role = "ADMIN"
RoleMember Role = "MEMBER"
)

// User represents an authenticated entity in the system
type User struct {
ID        string    `json:"id" dynamodbav:"pk"`
Email     string    `json:"email" dynamodbav:"email"`
Role      Role      `json:"role" dynamodbav:"role"`
Status    string    `json:"status" dynamodbav:"status"`
CreatedAt time.Time `json:"created_at" dynamodbav:"created_at"`
}

var (
ErrInvalidEmailDomain = errors.New("email domain not authorized")
)


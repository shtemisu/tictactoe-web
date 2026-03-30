package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Login    string
	Password string
}

type SignUpRequest struct {
	Login    string
	Password string
}

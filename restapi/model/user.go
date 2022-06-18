package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id                  uuid.UUID
	Username            string
	Email               string
	EmailVerified       bool
	PhoneNumber         string
	PhoneNumberVerified bool
	DisplayName         string
	Avatar              string
	Language            string
	Deleted             bool
	JoinedAt            time.Time
}

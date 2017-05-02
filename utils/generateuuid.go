package utils

import uuid "github.com/satori/go.uuid"

// GenerateUUID4 generates a UUID Version 4 based on RFC 4122
func GenerateUUID4() string {
	u4 := uuid.NewV4()
	return u4.String()
}

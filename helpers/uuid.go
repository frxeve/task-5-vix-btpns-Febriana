package helpers

import "github.com/google/uuid"

func GetUUID() string {
	uid := uuid.New()
	return uid.String()
}

package entity

import (
	"encoding/json"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func PasswordHash(password string) (string, error) {
	psw := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(psw, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func GetIdFromInsertedObjectId(userId primitive.ObjectID) string {
	newUserId := strings.TrimSuffix(strings.TrimPrefix(userId.String(), "ObjectID(\""), "\")")

	return newUserId
}

func MarshalIndentHelper(s interface{}) []byte {

	byteInput, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return nil
	}
	return byteInput
}

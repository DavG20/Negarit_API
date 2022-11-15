package entity

import (
	"math/rand"
	"strings"
	"time"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const Characters = "abcdefghijelmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateRandomString() string {
	Lists := make([]byte, 6)
	for i := range Lists {
		Lists[i] = Characters[seededRand.Intn(len(Characters))]
	}
	return string(Lists)
}

func Getextension(fileName string) string {
	fileNames := strings.Split(fileName, ".")
	if len(fileNames) > 1 {
		return fileNames[len(fileNames)-1]
	}
	return ""
}

func IsProfileValid(fileName string) bool {
	extension := Getextension(fileName)
	extension = strings.ToLower(extension)
	for _, e := range ImageExtensions {

		if e == extension {
			return true
		}
	}
	return false
}

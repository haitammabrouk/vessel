package id

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
)

const shortLen = 12

// GetShortHandId keeps a shorthand version of container id, Ideally grab the first 12 characters
func GetShortHandId(id string) string{
	trimTo := shortLen
	if len(id) < shortLen {
		trimTo = len(id)
	}
	return id[:trimTo]
}

func GenerateRandomId() string {
	b := make([]byte, 32)
	r := rand.Reader

	for {
		if _, err := r.Read(b); err != nil {
			panic(err)
		}
		id := hex.EncodeToString(b)
		// if the shorthand id can be parsed as int in base 10 then regenerate id
		// generated id is picked as default hostname for the container
		if _, err := strconv.ParseInt(GetShortHandId(id), 10, 64); err == nil {
			continue
		}
		return id
	}
}
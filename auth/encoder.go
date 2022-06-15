package auth

import (
	b64 "encoding/base64"
	"fmt"
)

func DecodeBase64ToString(data []byte) []byte {
	decoded, err := b64.URLEncoding.DecodeString(string(data))
	if err != nil {
		fmt.Println(err)
		// return ""
	}
	return decoded
}

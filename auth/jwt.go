package auth

import (
	localStructs "example/web-service-gin/structs"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GetMyRecordedAccessTokenLocal(email string) localStructs.PublicAccess {
	var recordFound localStructs.PublicAccess
	for i := 0; i < len(localStructs.PublicAccessTokenRecords); i++ {
		if localStructs.PublicAccessTokenRecords[i].Email == email {
			recordFound = localStructs.PublicAccessTokenRecords[i]
		}
	}
	return recordFound
}

func GetMyRecordedPrivateAccessLocal(email string, password string) localStructs.PrivateAccess {
	var recordFound localStructs.PrivateAccess
	for i := 0; i < len(localStructs.PrivateAccessTokenRecords); i++ {
		if localStructs.PrivateAccessTokenRecords[i].Email == email && localStructs.PrivateAccessTokenRecords[i].Password == password {
			recordFound = localStructs.PrivateAccessTokenRecords[i]
		}
	}
	return recordFound
}

func GenerateJWT(email, name string) (string, error) {
	var mySigningKey = []byte("YUSUF ADISAPUTRO")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["name"] = name
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func DecodeJWT(jwtToken string) *jwt.Token {
	var mySigningKey = []byte("YUSUF ADISAPUTRO")

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return mySigningKey, nil
	})

	if err != nil {
		return token
	}
	return token
}

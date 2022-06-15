package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	auth "example/web-service-gin/auth"
	baskerHandler "example/web-service-gin/pantry"
	localStructs "example/web-service-gin/structs"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type SocketEvent struct {
	SocketEvent string `json:"socketEvent"`
}
type ActiveSocket struct {
	SocketId string `json:"socketId"`
	// SocketRoom string `json:"socketRoom"`
	SocketEvents []SocketEvent `json:"socketEvents"`
}

var ActiveSockets = []ActiveSocket{}

func main() {
	router := gin.Default()

	router.POST("/request/access/:level/:access", requestPublicAccessToken)
	router.GET("/find/my/access", getMyRecordedAccessToken)
	router.POST("/decode/my/access", decodeAccessToken)
	router.POST("/remove/my/access", removeAccessToken)

	router.POST("/upload-data", uploadGlobalDataToBasket)
	router.GET("/download-data/:basket", getData)
	router.Run("localhost:3000")
}

func getMyRecordedAccessToken(c *gin.Context) {
	// email := c.Request.URL.Query().Get("email")
	ip := c.ClientIP()
	for i := 0; i < len(localStructs.PublicAccessTokenRecords); i++ {
		if localStructs.PublicAccessTokenRecords[i].IP == ip {
			c.IndentedJSON(http.StatusOK, localStructs.PublicAccessTokenRecords[i])
			return
		}
	}
	for i := 0; i < len(localStructs.PrivateAccessTokenRecords); i++ {
		if localStructs.PrivateAccessTokenRecords[i].IP == ip {
			c.IndentedJSON(http.StatusOK, localStructs.PrivateAccessTokenRecords[i])
			return
		}
	}
	c.IndentedJSON(http.StatusBadRequest, nil)
	return
}

func requestPublicAccessToken(c *gin.Context) {
	var newPublicRequestAccess struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	// Call BindJSON to bind the received JSON to
	// newPublicAccess.
	if err := c.BindJSON(&newPublicRequestAccess); err != nil {
		return
	}
	if len(newPublicRequestAccess.Password) <= 0 && c.Param("access") == "private" {
		c.IndentedJSON(http.StatusBadRequest, "Please provide password")
		return
	}
	// fmt.Println(auth.GetMyRecordedAccessTokenLocal(newPublicRequestAccess.Email) == publicAccess{})
	if auth.GetMyRecordedAccessTokenLocal(newPublicRequestAccess.Email) != (localStructs.PublicAccess{}) && c.Param("access") == "public" {
		var requestedAccessLevel localStructs.AccessLevel
		switch level := c.Param("level"); level {
		case "level1":
			c.IndentedJSON(http.StatusBadRequest, errorDefaultResponse(false, "You are trying to request a different access level while your param says public access request"))
			return
		case "level2":
			c.IndentedJSON(http.StatusBadRequest, errorDefaultResponse(false, "You are trying to request a different access level while your param says public access request"))
			return
		case "level3":
			c.IndentedJSON(http.StatusBadRequest, errorDefaultResponse(false, "You are trying to request a different access level while your param says public access request"))
			return
		case "Global":
			requestedAccessLevel = localStructs.Global
		default:
			c.IndentedJSON(http.StatusBadRequest, errorDefaultResponse(false, "You are trying to request a different access level while your param says public access request"))
			return
		}
		if auth.GetMyRecordedAccessTokenLocal(newPublicRequestAccess.Email).LevelAccess == requestedAccessLevel {
			c.IndentedJSON(http.StatusOK, auth.GetMyRecordedAccessTokenLocal(newPublicRequestAccess.Email))
		} else {
			c.IndentedJSON(http.StatusBadRequest, errorDefaultResponse(false, "You are trying to request a different access level while your old token is still active"))
		}
		return
	} else if auth.GetMyRecordedPrivateAccessLocal(newPublicRequestAccess.Email, newPublicRequestAccess.Password) != (localStructs.PrivateAccess{}) && c.Param("access") == "private" {
		var requestedAccessLevel localStructs.AccessLevel
		switch level := c.Param("level"); level {
		case "level1":
			requestedAccessLevel = localStructs.Level1
		case "level2":
			requestedAccessLevel = localStructs.Level2
		case "level3":
			requestedAccessLevel = localStructs.Level3
		case "Global":
			c.IndentedJSON(http.StatusBadRequest, errorDefaultResponse(false, "You are trying to request a Global access level while your param says private access request"))
			return
		default:
			c.IndentedJSON(http.StatusBadRequest, errorDefaultResponse(false, "You are trying to request a different access level while your param says private access request"))
			return
		}
		if auth.GetMyRecordedPrivateAccessLocal(newPublicRequestAccess.Email, newPublicRequestAccess.Password).LevelAccess == requestedAccessLevel {
			c.IndentedJSON(http.StatusOK, auth.GetMyRecordedPrivateAccessLocal(newPublicRequestAccess.Email, newPublicRequestAccess.Password))
		} else {
			c.IndentedJSON(http.StatusBadRequest, errorDefaultResponse(false, "You are trying to request a different access level while your old token is still active"))
		}
		return
	}

	validToken, err := auth.GenerateJWT(newPublicRequestAccess.Email, newPublicRequestAccess.Name)
	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	var accessLevel localStructs.AccessLevel
	level := c.Param("level")
	if level == "level1" && c.Param("access") == "private" {
		accessLevel = localStructs.Level1
		var newPrivateAccess localStructs.PrivateAccess
		newPrivateAccess.JWTtoken = validToken
		newPrivateAccess.Email = newPublicRequestAccess.Email
		newPrivateAccess.IP = c.ClientIP()
		newPrivateAccess.Time = time.Now().Add(time.Minute * 30).Unix()
		newPrivateAccess.LevelAccess = accessLevel
		newPrivateAccess.Password = newPublicRequestAccess.Password

		// Add the new album to the slice.
		localStructs.PrivateAccessTokenRecords = append(localStructs.PrivateAccessTokenRecords, newPrivateAccess)
		fmt.Println(auth.AlterBasketPrivateData(localStructs.PrivateAccessTokenRecords, "authPrivateRecords"))
		c.IndentedJSON(http.StatusCreated, newPrivateAccess)
		return
	} else if level == "Global" && c.Param("access") == "public" {
		accessLevel = localStructs.Global
		var newPublicAccess localStructs.PublicAccess
		newPublicAccess.JWTtoken = validToken
		newPublicAccess.Email = newPublicRequestAccess.Email
		newPublicAccess.IP = c.ClientIP()
		newPublicAccess.Time = time.Now().Add(time.Minute * 30).Unix()
		newPublicAccess.LevelAccess = accessLevel

		// Add the new album to the slice.
		localStructs.PublicAccessTokenRecords = append(localStructs.PublicAccessTokenRecords, newPublicAccess)
		fmt.Println(auth.AlterBasketData(localStructs.PublicAccessTokenRecords, "authRecords"))
		c.IndentedJSON(http.StatusCreated, newPublicAccess)
		return
	} else {
		c.IndentedJSON(http.StatusBadRequest, "no access available")
		return
	}
}

func decodeAccessToken(c *gin.Context) {
	if c.Request.Header["Token"] == nil {
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	token := auth.DecodeJWT(c.Request.Header["Token"][0])

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var recordFound interface{}
		for i := 0; i < len(localStructs.PublicAccessTokenRecords); i++ {
			if localStructs.PublicAccessTokenRecords[i].Email == claims["email"] {
				recordFound = localStructs.PublicAccessTokenRecords[i]
			}
		}
		for i := 0; i < len(localStructs.PrivateAccessTokenRecords); i++ {
			if localStructs.PrivateAccessTokenRecords[i].Email == claims["email"] {
				recordFound = localStructs.PrivateAccessTokenRecords[i]
			}
		}
		c.IndentedJSON(http.StatusOK, recordFound)
		return
	}
	c.IndentedJSON(http.StatusBadRequest, errorDefaultResponse(token.Valid, "JWT token is not available in records"))
}

func removeAccessToken(c *gin.Context) {
	if c.Request.Header["Token"] == nil {
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	var mySigningKey = []byte("YUSUF ADISAPUTRO")

	token, err := jwt.Parse(c.Request.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return mySigningKey, nil
	})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, errorDefaultResponse(token.Valid, "JWT token is invalid"))
		return
	}

	if token.Valid {
		var recordFound interface{}
		for i := 0; i < len(localStructs.PublicAccessTokenRecords); i++ {
			if localStructs.PublicAccessTokenRecords[i].JWTtoken == c.Request.Header["Token"][0] {
				localStructs.PublicAccessTokenRecords = append(localStructs.PublicAccessTokenRecords[:i], localStructs.PublicAccessTokenRecords[i+1:]...)
				recordFound = localStructs.PublicAccessTokenRecords[i].JWTtoken
			}
		}
		for i := 0; i < len(localStructs.PrivateAccessTokenRecords); i++ {
			if localStructs.PrivateAccessTokenRecords[i].JWTtoken == c.Request.Header["Token"][0] {
				localStructs.PrivateAccessTokenRecords = append(localStructs.PrivateAccessTokenRecords[:i], localStructs.PrivateAccessTokenRecords[i+1:]...)
				recordFound = localStructs.PublicAccessTokenRecords[i].JWTtoken
			}
		}
		c.IndentedJSON(http.StatusGone, recordFound)
		return
	}
	c.IndentedJSON(http.StatusBadRequest, errorDefaultResponse(token.Valid, "JWT token is not available in records"))
}

func errorDefaultResponse(status bool, reason string) localStructs.DefaultErrorMessage {
	var err localStructs.DefaultErrorMessage
	err.Status = status
	err.Reason = reason
	return err
}

func uploadGlobalDataToBasket(c *gin.Context) {
	theFile, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	// Get raw file bytes - no reader method
	openedFile, _ := theFile.Open()
	thefileData, _ := ioutil.ReadAll(openedFile)

	if theFile.Size >= 100000 {
		c.IndentedJSON(http.StatusBadRequest, theFile.Size)
		return
	}
	mimeType := http.DetectContentType(thefileData)
	var theStructFile baskerHandler.File
	theStructFile.FileData = thefileData
	theStructFile.MimeType = mimeType

	c.IndentedJSON(http.StatusAlreadyReported, baskerHandler.AppendDataInBasket(theStructFile, "files"))

}

func getData(c *gin.Context) {
	c.IndentedJSON(http.StatusFound, baskerHandler.DecodeBasketData(c.Param("basket")))
}

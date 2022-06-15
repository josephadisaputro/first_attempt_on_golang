package baskerHandler

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	auth "example/web-service-gin/auth"
)

type File struct {
	FileData []byte `json:"fileData"`
	MimeType string `json:"mimeType"`
}

func AlterBasketData(payloadRequest []File, basket string) bool {
	url := "https://getpantry.cloud/apiv1/pantry/89220f30-5493-47cb-a8c7-867bdb30ea34/basket/" + basket
	method := "POST"

	var jsonnize struct {
		Data []File `json:"data"`
	}
	jsonnize.Data = payloadRequest

	jsonizeReorderPayload, err := json.Marshal(jsonnize)
	if err != nil {
		fmt.Println(err)
		return false
	}

	payload := strings.NewReader("{\"data\": \"" + b64.URLEncoding.EncodeToString([]byte(string(jsonizeReorderPayload))) + "\"}")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return false
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(string(body))
	return true
}

func AppendDataInBasket(data File, basket string) bool {
	var makeStruct struct {
		Data []File `json:"data"`
	}
	_ = json.Unmarshal([]byte(auth.GetBasketData(basket)), &makeStruct)

	for i := 0; i < len(makeStruct.Data); i++ {
		if string(makeStruct.Data[i].FileData) == string(data.FileData) {
			return false
		}
	}

	makeStruct.Data = append(makeStruct.Data, data)

	return AlterBasketData(makeStruct.Data, "files")
}

func DecodeBasketData(basket string) []File {
	var makeStruct struct {
		Data []File `json:"data"`
	}
	_ = json.Unmarshal([]byte(auth.GetBasketData(basket)), &makeStruct)
	return makeStruct.Data
}

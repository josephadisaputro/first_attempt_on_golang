package auth

import (
	b64 "encoding/base64"
	"encoding/json"
	localStructs "example/web-service-gin/structs"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func AlterBasketData(payloadRequest []localStructs.PublicAccess, basket string) interface{} {
	url := "https://getpantry.cloud/apiv1/pantry/89220f30-5493-47cb-a8c7-867bdb30ea34/basket/" + basket
	method := "POST"

	var jsonnize struct {
		JWTrecords []localStructs.PublicAccess `json:"jwtRecords"`
	}
	jsonnize.JWTrecords = payloadRequest

	jsonizeReorderPayload, err := json.Marshal(jsonnize)
	if err != nil {
		fmt.Println(err)
		return err
	}

	payload := strings.NewReader("{\"data\": \"" + b64.URLEncoding.EncodeToString([]byte(string(jsonizeReorderPayload))) + "\"}")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// fmt.Println(string(body))
	return string(body)
}

func AlterBasketPrivateData(payloadRequest []localStructs.PrivateAccess, basket string) interface{} {
	url := "https://getpantry.cloud/apiv1/pantry/89220f30-5493-47cb-a8c7-867bdb30ea34/basket/" + basket
	method := "POST"

	var jsonnize struct {
		JWTrecords []localStructs.PrivateAccess `json:"jwtRecords"`
	}
	jsonnize.JWTrecords = payloadRequest

	jsonizeReorderPayload, err := json.Marshal(jsonnize)
	if err != nil {
		fmt.Println(err)
		return err
	}

	payload := strings.NewReader("{\"data\": \"" + b64.URLEncoding.EncodeToString([]byte(string(jsonizeReorderPayload))) + "\"}")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// fmt.Println(string(body))
	return string(body)
}

func GetBasketData(basket string) string {
	url := "https://getpantry.cloud/apiv1/pantry/89220f30-5493-47cb-a8c7-867bdb30ea34/basket/" + basket
	method := "GET"

	payload := strings.NewReader(string(""))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var responseJson struct {
		Data string `json:"data"`
	}
	json.Unmarshal([]byte(string(body)), &responseJson)
	var decodedByte, _ = b64.URLEncoding.DecodeString(responseJson.Data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(decodedByte)
}

// func CompareBasketDataWithJWTLocalStructs(basket string, accessList []localStructs.PublicAccess) string {
// 	basketDate := GetBasketData(basket)
// 	return basketDate
// }

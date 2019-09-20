package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"bytes"
	"os"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

const PROJECT_ID = "remote-config-example-b25bb"
const BASE_URL = "https://firebaseremoteconfig.googleapis.com"
const REMOTE_CONFIG_ENDPOINT = "v1/projects/" + PROJECT_ID + "/remoteConfig"
const REMOTE_CONFIG_URL = BASE_URL + "/" + REMOTE_CONFIG_ENDPOINT

var client = &http.Client{}

func serviceAccount(credentialFile string) (*oauth2.Token, error) {
	b, err := ioutil.ReadFile(credentialFile)
	if err != nil {
		return nil, err
	}
	var c = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}
	json.Unmarshal(b, &c)
	config := &jwt.Config{
		Email:      c.Email,
		PrivateKey: []byte(c.PrivateKey),
		Scopes: []string{
			"https://www.googleapis.com/auth/firebase.remoteconfig",
		},
		TokenURL: google.JWTTokenURL,
	}
	token, err := config.TokenSource(oauth2.NoContext).Token()
	if err != nil {
		return nil, err
	}
	return token, nil
}

func publish(token *oauth2.Token, Etag string) string {

	// read config.json
	readFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println(err)
	}

	request, err := http.NewRequest("PUT", REMOTE_CONFIG_URL, bytes.NewReader(readFile))
	if err != nil {
		log.Fatal("Error : %v\n", err)
	}

	request.Header.Set("Authorization", "Bearer "+token.AccessToken)
	request.Header.Add("Content-Type", "application/json; UTF-8")
	request.Header.Add("If-Match", Etag)

	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	// if resp.Status is 200
	if response.StatusCode == http.StatusOK {

		// Read response body
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Convert Response Body to String
		bodyString := string(bodyBytes)

		// Print it out
		fmt.Println(bodyString)
		
		// Get ETag
		fmt.Printf("Found ETag: %+v\n", response.Header["Etag"][0])

		return string(response.Header["Etag"][0])
	}
	return ""
}

func rollbackVersion(token *oauth2.Token, version int) {
	
	// client := &http.Client{}

	respJSON := map[string]int{"version_number": version}
	b, _ := json.Marshal(respJSON)

	// Create new request
	req, err := http.NewRequest(http.MethodPost, REMOTE_CONFIG_URL+":rollback", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	if resp.StatusCode == http.StatusOK {
		// Read response body
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		// Convert Response Body to String
		bodyString := string(bodyBytes)
		print(bodyString)

		// get latest etag
		etag := getRemoteConfig(token)
		writeEtag(etag)
	}

}

func getRemoteConfig(token *oauth2.Token) (string) {

	req, err := http.NewRequest("GET", REMOTE_CONFIG_URL, nil)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	// Set Authorization Header
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	// if resp.Status is 200
	if resp.StatusCode == http.StatusOK {

		// Read response body
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		// Convert Response Body to String
		bodyString := string(bodyBytes)
		// Print it out
		fmt.Println(bodyString)
		// Get ETag
		fmt.Printf("Found ETag: %+v\n", resp.Header["Etag"][0])

		// Write the response body into config.json
		_ = ioutil.WriteFile("config.json", bodyBytes, 0644)

		return string(resp.Header["Etag"][0])
	}
	return ""
}

func listVersion(token *oauth2.Token, size int) {

	req, err := http.NewRequest(http.MethodGet, REMOTE_CONFIG_URL+":listVersions?pageSize="+strconv.Itoa(size), nil)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	// Set Authorization Header
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	if resp.StatusCode == http.StatusOK {
		// Read response body
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		// Convert Response Body to String
		bodyString := string(bodyBytes)
		print("%+v", resp)
		print(bodyString)
	}
	
}

func readEtag() string {
	data, err := ioutil.ReadFile("etag.txt")
    if err != nil {
        fmt.Println("File reading error", err)
    }
	etag := string(data)
	return etag
}

func writeEtag(etag string) {
	f, err := os.Create("etag.txt")
    if err != nil {
        fmt.Println(err)
        return
	}

	f.WriteString(etag)
	f.Close()
}

func main() {

	fmt.Println("Firebase Remote-Config Console\n")
	fmt.Println("1. Get Remote Config")
	fmt.Println("2. Publish")
	fmt.Println("3. Rollback")
	fmt.Println("4. Show Version List")
	fmt.Print("\nOperation: ")
	var i int
    fmt.Scanf("%d", &i)

	// get access token
	token, err := serviceAccount("service-account.json")
	if err != nil {
		fmt.Println("Error acquiring token: %v", err)
	}
	
	if i == 1 {
		writeEtag(getRemoteConfig(token))
	} 
	if i == 2 {
		writeEtag(publish(token, readEtag()))
	}
	if i == 3 {
		fmt.Print("Rollback to version : ")
		var version int
		fmt.Scanf("%d", &version)
		rollbackVersion(token, version)
	}
	if i == 4 {
		fmt.Print("\nSize list : ")
		var size int
		fmt.Scanf("%d", &size)
		listVersion(token, size)
	}

}
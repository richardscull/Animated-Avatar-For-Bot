package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var token string
var discordAPI = "https://discord.com/api/v10"

func main() {
	fmt.Println("Please enter your bot's token:")
	fmt.Scanln(&token)

	file, err := os.Open("avatar.gif")
	if err != nil {
		fmt.Println("Error opening file:", err)
		fmt.Println("Please make sure you placed your wanted avatar in the same directory as the executable.")
		fmt.Println("The file should be named 'avatar.gif' and should be a .gif file.")
		return
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	encodedImage := base64.StdEncoding.EncodeToString(imageData)

	payload, err := json.Marshal(map[string]interface{}{
		"username": getUsername(),
		"avatar":   "data:image/gif;base64," + encodedImage,
	})
	if err != nil {
		fmt.Println("Error creating request body:", err)
		return
	}

	url := discordAPI + "/users/@me"
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bot "+token)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error response status:", resp.Status)
		fmt.Println("Check if your token is correct and if your bot has the necessary permissions.")
		return
	}

	fmt.Println("Avatar changed successfully! It may take a few seconds to update.")
	fmt.Scanf("h")
}

func getUsername() string {
	url := discordAPI + "/users/@me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Authorization", "Bot "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error response status:", resp.Status)
		fmt.Println("Check if your token is correct and if your bot has the necessary permissions.")
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	var response struct {
		Username string `json:"username"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error parsing response body:", err)
		return ""
	}

	return response.Username
}

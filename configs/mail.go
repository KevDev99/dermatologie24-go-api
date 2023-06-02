package configs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func SendMail(payload string) {
	url := os.Getenv("BREVO_URL") + "/smtp/email"
	apiKey := os.Getenv("BREVO_API_KEY")

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err.Error())
		return
	}

	fmt.Println(string(body))

	fmt.Println("Response status:", resp.Status)
}

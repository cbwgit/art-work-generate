package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type ArtworkOrder struct {
	ID     string `json:"orderId"`
	URL    string `json:"url"`
	Prompt string `json:"prompt"`
}

func generateArt(prompt string, authorization string) (*ArtworkOrder, error) {
	url := "https://api.neural.love/v1/ai-art/generate"

	payload := strings.NewReader(fmt.Sprintf("{\"style\":\"painting\",\"layout\":\"square\",\"amount\":4,\"isPublic\":true,\"isPriority\":false,\"isHd\":false,\"steps\":25,\"cfgScale\":7.5,\"prompt\":\"%s\"}", prompt))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Bearer "+authorization)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var artworkOrder ArtworkOrder
	err = json.Unmarshal(body, &artworkOrder)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Create art order: %+v\n", artworkOrder)

	return &artworkOrder, nil
}

func main() {
	// Set Gin framework to output logs to the terminal
	gin.DefaultWriter = os.Stdout
	router := gin.Default()
	router.LoadHTMLGlob("templates/*") // Load templates

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.POST("/submit", func(c *gin.Context) {
		prompt := c.PostForm("prompt") // Get the prompt value from the user submission

		authorization := "" // Initialize with an empty authorization code

		// Create artwork order
		artworkOrder, err := generateArt(prompt, authorization)
		if err != nil {
			fmt.Println("Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error"})
			return
		}

		// Pass order information to the index template
		c.HTML(http.StatusOK, "index.html", gin.H{
			"orderID":  artworkOrder.ID,
			"orderURL": artworkOrder.URL,
			"prompt":   artworkOrder.Prompt,
		})
	})

	router.POST("/authorize", func(c *gin.Context) {
		authorization := c.PostForm("authorization") // Get the authorization value from the user submission

		// Handle the authorization logic here...
		fmt.Println("Authorization:", authorization)

		// Pass the authorization code to the index template
		c.HTML(http.StatusOK, "index.html", gin.H{
			"authorization": authorization,
		})
	})

	// Set the port
	address := fmt.Sprintf(":%d", 8080) // Use port 8080

	router.Run(address)
}

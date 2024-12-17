package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	baseURL := "http://localhost:8080/download/large"

	for i := 1; i <= 14; i++ {
		// Construct URL with increasing file number
		url := baseURL + strconv.Itoa(i) + ".dmg"
		filename := "newfile" + strconv.Itoa(i) + ".dmg"

		// Download the file
		fmt.Println(url)
		err := downloadFile(url, filename)
		if err != nil {
			fmt.Printf("Failed to download file: %s\n", err)
			return
		}
		fmt.Printf("Downloaded '%s' successfully.\n", filename)
	}
}

// downloadFile downloads a file from the given URL and saves it with the specified filename.
func downloadFile(url, filename string) error {
	// Create the file
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the HTTP request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

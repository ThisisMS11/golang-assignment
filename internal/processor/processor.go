package processor

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

// ProcessImage downloads and processes an image, returning its perimeter
func ProcessImage(imageURL string) (float64, error) {
	// Parse URL
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return 0, fmt.Errorf("invalid URL: %v", err)
	}

	// Make HTTP request
	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return 0, fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to download image: HTTP %d", resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to decode image: %v", err)
	}

	// Calculate perimeter
	bounds := img.Bounds()
	height := bounds.Max.Y - bounds.Min.Y
	width := bounds.Max.X - bounds.Min.X
	perimeter := 2.0 * float64(height+width)

	// simulating GPU
	sleepTime := 100 + rand.Intn(300)
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	return perimeter, nil
}

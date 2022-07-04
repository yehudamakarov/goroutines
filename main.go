package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	//
}

func getImageFromEbay() {
	res, err := http.Get("https://i.ebayimg.com/thumbs/images/g/j1wAAOSwujZiUMBu/s-l225.jpg")
	if err != nil {
		log.Fatal(fmt.Errorf("couldn't get image from ebay: %w", err))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(fmt.Errorf("couldn't close: %w", err))
		}
	}(res.Body)

	// read
	readAll, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(fmt.Errorf("couldn't read bytes of response %w", err))
	}

	err = os.WriteFile("./result.jpg", readAll, 0644)
	if err != nil {
		log.Fatal(fmt.Errorf("couldn't write file: %w", err))
	}
}

func downloadAndUploadImages(items []int) chan string {
	lengthOfItems := len(items)
	images := goGetImages(items, lengthOfItems)
	imageUrls := goUploadImages(images, lengthOfItems)
	return imageUrls
}

func goGetImages(items []int, lengthOfItems int) chan string {
	images := make(chan string, lengthOfItems)

	var wg sync.WaitGroup
	wg.Add(lengthOfItems)

	go getImages(items, images, &wg)

	return images
}

func getImages(items []int, images chan string, wg *sync.WaitGroup) {
	for _, item := range items {
		go callApiForImage(item, images, wg)
	}

	wg.Wait()
	close(images)
}

func callApiForImage(item int, images chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	image := getImage(item)
	images <- image
}

func getImage(item int) string {
	image := fmt.Sprintf("image: %d", item)
	fmt.Printf("getting %s...\n", image)
	randomSleep()
	return image
}

func randomSleep() {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(3)
	time.Sleep(time.Duration(n) * time.Second)
}

func goUploadImages(images chan string, lengthOfItems int) chan string {
	imageUrls := make(chan string, lengthOfItems)
	go uploadImages(images, imageUrls, lengthOfItems)
	return imageUrls
}

func uploadImages(
	images chan string, imageUrls chan string, lengthOfItems int,
) {
	var wg sync.WaitGroup
	wg.Add(lengthOfItems)

	for image := range images {
		go uploadImage(image, imageUrls, &wg)
	}

	wg.Wait()
	close(imageUrls)
}

func uploadImage(image string, imageUrls chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	imageUrl := storeImage(image)
	imageUrls <- imageUrl
}

func storeImage(image string) string {
	imageUrl := fmt.Sprintf("URL for %s", image)
	fmt.Printf("storing %s...\n", imageUrl)
	randomSleep()
	return imageUrl
}

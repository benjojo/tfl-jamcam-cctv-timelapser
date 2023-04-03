package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	flag.Parse()
	if *honkServer == "" || *honkUsername == "" || *honkPassword == "" {
		flag.Usage()
		os.Exit(1)
	}
	getHonkSiteToken() // verify the creds

	log.Printf("Fetching TFL XML")
	req, _ := http.NewRequest("GET", "https://content.tfl.gov.uk/camera-list.xml", nil)
	req.Header.Set("User-Agent", "tfl-api@benjojo.co.uk")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to request XML: %v", err)
	}
	xmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to download XML: %v", err)
	}

	cam1 := getCameraFromXML(xmlBytes)
	started := time.Now()
	minTicker := time.NewTicker(time.Minute * 5)
	images := make([][]byte, 0)
	for {
		if time.Since(started) > (time.Hour * 24) {
			break
		}
		<-minTicker.C

		req, _ := http.NewRequest("GET", fmt.Sprintf("https://s3-eu-west-1.amazonaws.com/jamcams.tfl.gov.uk/%s.jpg?%d",
			cam1.ID, time.Now().Unix()), nil)

		req.Header.Set("User-Agent", "tfl-api@benjojo.co.uk")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Failed to get image: %v", err)
			continue
		}

		imageBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to get image: %v", err)
			continue
		}

		images = append(images, imageBytes)
	}

	videoFile := createTimeLapseVideo(images)
	tok := getHonkSiteToken()
	postToHonk(tok, cam1.ID, videoFile)
}

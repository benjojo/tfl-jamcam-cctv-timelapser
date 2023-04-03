package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

var (
	honkServer   = flag.String("server", "benjojo.co.uk", "server to connnect")
	honkUsername = flag.String("username", "tfltimelapse", "username to use")
	honkPassword = flag.String("password", "", "password to use")
)

func getHonkSiteToken() string {
	form := make(url.Values)
	form.Add("username", *honkUsername)
	form.Add("password", *honkPassword)
	form.Add("gettoken", "1")
	loginurl := fmt.Sprintf("https://%s/dologin", *honkServer)
	req, err := http.NewRequest("POST", loginurl, strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	answer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("status: %d: %s", resp.StatusCode, answer)
	}
	return string(answer)
}

func postToHonk(authToken string, camName string, videoFile []byte) {
	buf := bytes.NewBuffer(nil)
	formWriter := multipart.NewWriter(buf)

	// multipart.Form.
	w, _ := formWriter.CreateFormField("token")
	w.Write([]byte(authToken))
	w, _ = formWriter.CreateFormField("action")
	w.Write([]byte("honk"))
	w, _ = formWriter.CreateFormField("noise")
	w.Write([]byte(fmt.Sprintf("tfltimelapse presents: %v ", camName)))
	w, _ = formWriter.CreateFormFile("donk", "vid.mp4")
	io.Copy(w, bytes.NewReader(videoFile))
	formWriter.Close()

	apiurl := fmt.Sprintf("https://%s/api", "benjojo.co.uk")
	req, err := http.NewRequest("POST", apiurl, buf)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", formWriter.FormDataContentType())
	req.Header.Add("Authorization", authToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	answer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("status: %d: %s", resp.StatusCode, answer)
	}
}

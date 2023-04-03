package main

import (
	"bytes"
	"io"
	"os/exec"
)

func createTimeLapseVideo(jpegs [][]byte) (video []byte) {
	ffmpegPid := exec.Command("ffmpeg", "-f", "mjpeg", "-i", "pipe:0", "-movflags", "frag_keyframe+empty_moov", "-f", "mp4", "pipe:1")
	inputImages, _ := ffmpegPid.StdinPipe()
	outputVideo, _ := ffmpegPid.StdoutPipe()

	ffmpegPid.Start()

	videoBuf := bytes.Buffer{}
	outputClosed := make(chan bool)
	go func() {
		io.Copy(&videoBuf, outputVideo)
		outputClosed <- true
		ffmpegPid.Wait()
	}()

	// f, _ := os.Create("debug.mjpeg")
	for _, img := range jpegs {
		fin := []byte("\n--myboundary\nContent-Type: image/jpeg\n\n")
		fin = append(fin, img...)
		inputImages.Write(fin)
		// f.Write(fin)
	}

	inputImages.Close()
	<-outputClosed

	// MP4 should be in videoBuf
	return videoBuf.Bytes()
}

package main

import "testing"

func TestHonkPost(t *testing.T) {
	t.SkipNow()
	jpegs := makeTestJpegDataset()
	video := createTimeLapseVideo(jpegs)
	tok := getHonkSiteToken()
	postToHonk(tok, "gotest", video)
}

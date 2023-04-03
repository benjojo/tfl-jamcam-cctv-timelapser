package main

import (
	"encoding/xml"
	"log"
	"math/rand"
	"time"
)

func getCameraFromXML(b []byte) XMLCam {
	DecodedXML := XMLFeed{}
	err := xml.Unmarshal(b, &DecodedXML)
	if err != nil {
		log.Fatalf("Unable to parse TFL XML: %v", err)
	}

	applicableEntries := make([]XMLCam, 0)
	for _, v := range DecodedXML.Cams.CameraList {
		if v.Available {
			applicableEntries = append(applicableEntries, v)
		}
	}

	rand.Seed(time.Now().UnixNano())

	if len(applicableEntries) != 0 {
		return applicableEntries[rand.Intn(len(applicableEntries)-1)]
	}

	return XMLCam{}
}

type XMLFeed struct {
	Cams XMLCamList `xml:"cameraList"`
}

type XMLCamList struct {
	CameraList []XMLCam `xml:"camera"`
	Rooturl    string   `xml:"rooturl"`
}

type XMLCam struct {
	/*
			<camera id="00001.01445" available="true">
		      <corridor/>
		      <location>A10 Grt Cambridge Rd/A110 S' bury Rd</location>
		      <currentView>A10 South</currentView>
		      <file>0000101445.jpg</file>
		      <captureTime>2012-04-20T15:15:03.046+01:00</captureTime>
		      <easting>534286.0</easting>
		      <northing>196318.0</northing>
		      <lat>51.649628</lat>
		      <lng>-0.060391203</lng>
		      <osgr>TQ343963</osgr>
		      <postCode>EN1 1RQ</postCode>
		    </camera>
	*/
	Location    string `xml:"location"`
	ID          string `xml:"id,attr"`
	Available   bool   `xml:"available,attr"`
	CurrentView string `xml:"currentView"`
	File        string `xml:"file"`
	CaptureTime string `xml:"captureTime"`
	Easting     string `xml:"easting"`
	Northing    string `xml:"northing"`
	Lat         string `xml:"lat"`
	Lng         string `xml:"lng"`
	Osgr        string `xml:"osgr"`
	PostCode    string `xml:"postCode"`
}

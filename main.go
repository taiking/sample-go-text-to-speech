package main

import (
	"context"
	"io/ioutil"
	"log"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/gin-gonic/gin"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)


type RequestJson struct {
	Name string
	Text string
}

func main() {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/create", func(c *gin.Context) {
		var json RequestJson
		c.BindJSON(&json)
		makeAudio(json.Name, json.Text)
		c.JSON(200, gin.H{
			"status": "success",
		})
	})
	r.Static("/sounds", "./sounds")
	r.Run()
}

func makeAudio(name string, text string) {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			Name: name,
			LanguageCode: "en-US",
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}

	filename := "sounds/output.mp3"
	err = ioutil.WriteFile(filename, resp.AudioContent, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
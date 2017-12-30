package main

import (
	"bufio"
	"bytes"
	"flag"
	"github.com/bwmarrin/discordgo"
	"github.com/jpoz/gomeme"
	"image/jpeg"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var imgPath string

func main() {
	token := flag.String("token", "", "token to use to connect to discord")
	flag.StringVar(&imgPath, "img", "", "path to image to use")

	flag.Parse()

	if *token == "" {
		log.Fatal("Please provide a token")
	}

	if imgPath == "" {
		log.Fatal("Please provide a path to the image to use")
	}

	dg, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatal(err)
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func ready(s *discordgo.Session, e *discordgo.Ready) {
	log.Println(e)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Respond to !slabbot
	if strings.HasPrefix(m.Content, "!slabbot") {
		// Split on space and find whatever is after the command
		idx := strings.Index(m.Content, " ")

		if idx < 0 {
			_, err := s.ChannelMessageSend(m.ChannelID, "i need to know what to put on the bottom line of the image, idiot")
			if err != nil {
				log.Println(err)
			}
			return
		}

		img, err := createMeme(m.Content[idx:])
		if err != nil {
			log.Println(err)
		}

		msg := discordgo.MessageSend{
			Files: []*discordgo.File{img},
		}

		_, err = s.ChannelMessageSendComplex(m.ChannelID, &msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func createMeme(bottomText string) (*discordgo.File, error) {
	f := &discordgo.File{}

	// Read image
	image, err := os.Open(imgPath)
	if err != nil {
		return f, err
	}

	defer image.Close()

	j, err := jpeg.Decode(image)
	if err != nil {
		return f, err
	}

	config := gomeme.NewConfig()

	config.FontSize = 80
	config.TopText = "Slab, did you just"
	config.BottomText = bottomText

	meme := &gomeme.Meme{
		Config:   config,
		Memeable: gomeme.JPEG{j},
	}

	b := bytes.Buffer{}
	r := bufio.NewWriter(&b)

	err = meme.Write(r)
	if err != nil {
		return f, err
	}

	f.Name = "slabdidyoujust.jpg"
	f.Reader = &b

	return f, nil
}

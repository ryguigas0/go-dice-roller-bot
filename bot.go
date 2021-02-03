package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("ERROR LOADING .env FILE")
	}
	Token = os.Getenv("TOKEN")
}

func main() {

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	if err = dg.Open(); err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if strings.Join(strings.Split(m.Content, "")[:2], "") != "&d" {
		return
	}

	num, err := strconv.Atoi(strings.Join(strings.Split(m.Content, "")[3:], ""))
	if err != nil || num <= 0 {
		s.ChannelMessageSend(m.ChannelID, "Por favor escreva um número (&d [número])")
		return
	}
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(num-1) + 1
	fmt.Printf("%s : %v : %v\n", m.Author, num, randNum)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%v tirou %v no d%v", m.Author.Mention(), randNum, num))
}

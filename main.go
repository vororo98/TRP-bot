package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	Token     string
	BotPrefix string

	config *configStruct
)

type configStruct struct {
	Token     string `json: "Token"`
	BotPrefix string `json: "BotPrefix"`
}

func ReadConfig() error {
	file, err := os.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(file))

	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix

	return nil
}

var BotId string
var goBot *discordgo.Session

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}
	if m.Content == BotPrefix+"ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	} else if strings.Contains(m.Content, "roll:") { //I want to parse a command like roll:1d6 or roll:10d20
		expRes := regexp.MustCompile("[0-9]+")
		result := expRes.FindAllStringSubmatch(m.Content, -1)

		rolls, err := strconv.Atoi(result[0][0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		rollType, err := strconv.Atoi(result[1][0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		message := "You rolled:"
		total := 0

		for i := 0; i < rolls; i++ {
			roll := rand.Intn(rollType + 1)
			total += roll
			message += " " + strconv.Itoa(roll) + ","
		}

		message += " = " + strconv.Itoa(total)
		_, _ = s.ChannelMessageSend(m.ChannelID, message)
	}
}

func main() {
	err := ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	Start()

	<-make(chan struct{})
	return
}

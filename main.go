package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Rompei/discord-chara-bot/bot"
)

func main() {
	configFileName := os.Getenv("CONFIG_FILE")
	if configFileName == "" {
		fmt.Println("config file does not exist.")
		return
	}
	config, err := bot.NewConfig(configFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	bots := make([]*bot.Bot, len(config.BotConfigs))
	for i := 0; i < len(config.BotConfigs); i++ {
		bots[i], err = bot.NewBot(&config.BotConfigs[i])
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	fmt.Println("Closing sessions.")
	for i := 0; i < len(bots); i++ {
		bots[i].Close()
	}
}

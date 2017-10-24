package bot

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Bot is bot structure.
type Bot struct {
	sounds map[string][][]byte
	dg     *discordgo.Session
	config *BotConfig
}

// NewBot is constructor.
func NewBot(config *BotConfig) (*Bot, error) {
	bot := &Bot{
		sounds: make(map[string][][]byte),
		config: config,
	}
	if err := bot.loadSounds(config.SoundDir); err != nil {
		return nil, err
	}

	dg, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		return nil, err
	}
	dg.AddHandler(bot.ready)
	dg.AddHandler(bot.messageCreate)
	dg.AddHandler(bot.guildCreate)
	bot.dg = dg
	err = bot.dg.Open()
	return bot, err
}

// Close the session.
func (b *Bot) Close() {
	b.dg.Close()
}

// Ready is called when bot is ready.
func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	s.UpdateStatus(0, b.config.BotPlaying)
}

// MessageCreate is called when message is created.
func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, "!") {
		return
	}

	if strings.HasPrefix(m.Content, "!"+b.config.BotPrefix+" ls") {
		res := "コマンド\n"
		for k := range b.sounds {
			res += "!" + k + "\n"
		}
		s.ChannelMessageSend(m.ChannelID, res)
		return
	}

	// Find the channel that the message came from.
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		// Could not find channel.
		return
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		// Could not find guild.
		return
	}

	// Look for the message sender in that guild's current voice states.
	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			err = b.playSound(s, g.ID, vs.ChannelID, strings.Replace(m.Content, "!", "", -1))
			if err != nil {
				fmt.Println("Error playing sound:", err)
			}
			return
		}
	}
}

// GuildCreate is called when guild joins.
func (b *Bot) guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, b.config.BotHello)
			return
		}
	}
}

func (b *Bot) loadSounds(dir string) (err error) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".dca") {
			cmdName := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			fmt.Println("loading", cmdName)
			if err = b.addSound(cmdName, path); err != nil {
				return err
			}
		}
		return nil
	})
	return
}

// loadSound attempts to load an encoded sound file from disk.
func (b *Bot) addSound(name, p string) (err error) {

	file, err := os.Open(p)
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return
	}

	var opuslen int16

	var sound [][]byte

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			b.sounds[name] = sound
			return err
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return
		}

		// Append encoded pcm data to the buffer.
		sound = append(sound, InBuf)
	}
	return
}

// playSound plays the current buffer to the provided channel.
func (b *Bot) playSound(s *discordgo.Session, guildID, channelID, cmdName string) (err error) {

	sound, ok := b.sounds[cmdName]
	if !ok {
		return errors.New("Not found sound")
	}

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range sound {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}

package bot

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Config is config for bots.
type Config struct {
	BotConfigs []BotConfig `yaml:"bots"`
}

// NewConfig is constructor.
func NewConfig(fname string) (*Config, error) {
	var config Config
	b, err := ioutil.ReadFile(fname)
	if err = yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// BotConfig is config for a bot.
type BotConfig struct {
	BotToken   string `yaml:"botToken"`
	BotPrefix  string `yaml:"botPrefix"`
	BotHello   string `yaml:"botHello"`
	BotPlaying string `yaml:"botPlaying"`
	SoundDir   string `yaml:"soundDir"`
}

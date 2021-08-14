package main

import (
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
	"github.com/gookit/validate"
)

// Config represents configuration
type Config struct {
	ListenAddr string `default:":80" validate:"required`

	// Mb
	BufferSize uint `default:"32" validate:"required|min:1`

	// Kb
	Bitrate uint `default:"32" validate:"required|min:8|max:320"`

	Codec       string `default:"libopus"   validate:"required"`
	Format      string `default:"opus"      validate:"required" `
	ContentType string `default:"audio/ogg" validate:"required" `
}

var config *Config

func initConfig(filePath string) error {
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		return err
	}

	if err := defaults.Set(config); err != nil {
		return err
	}

	if v := validate.Struct(config); !v.Validate() {
		return errors.New(v.Errors.One())
	}

	return nil
}

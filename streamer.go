package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/djherbis/nio/v3"
	"gopkg.in/djherbis/buffer.v1"
)

func stream(url string, w http.ResponseWriter) error {
	if url == "" {
		err := errors.New("empty url")

		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	w.Header().Set("Content-Type", config.ContentType)
	w.Header().Set("Accept-Ranges", "none")
	w.WriteHeader(http.StatusOK)

	return encode(url, w)
}

func encode(url string, w io.Writer) error {
	cmd := exec.Command("ffmpeg",
		"-i", url,
		"-vn",
		"-c:a", config.Codec,
		"-b:a", fmt.Sprintf("%dk", config.Bitrate),
		"-f", config.Format,
		"-ac", "2",
		"pipe:1")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer out.Close()

	err = cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Println(err)
		}
	}()

	buf := buffer.New(int64(config.BufferSize) * 1024 * 1024)
	_, err = nio.Copy(w, out, buf)
	if err != nil {
		cmd.Process.Kill()
	}

	return nil
}

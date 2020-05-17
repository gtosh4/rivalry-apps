package main

import (
	"errors"
	"flag"
	"io"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/gtosh4/rivalry-apps/internal/app"
	"github.com/sirupsen/logrus"
)

func main() {
	var token string
	flag.StringVar(&token, "token", "", "bot token")
	debug := false
	flag.BoolVar(&debug, "debug", false, "debug mode")

	flag.Parse()

	log := logrus.New()
	if debug {
		log.SetLevel(logrus.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	if token == "" {
		flag.Usage()
		log.Fatal("Token required")
		return
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %+v", err)
		return
	}

	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection to discord gateway: %+v", err)
		return
	}
	log.Info("Connected to discord gateway")

	srv := app.NewServer(log, dg, ":9004")

	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, io.EOF) {
		log.Errorf("Server exited due to error: %+v", err)
	}
}

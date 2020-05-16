package app

import (
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

type Server struct {
	Log logrus.FieldLogger

	srv     *http.Server
	discord *discordgo.Session
}

func NewServer(log logrus.FieldLogger, discord *discordgo.Session, addr string) *Server {
	s := &Server{
		Log:     log,
		discord: discord,
	}

	router := gin.New()
	router.Use(ginlogrus.Logger(log), gin.Recovery())

	router.POST("/:gid/googleApp", s.handleGoogleApp)

	s.srv = &http.Server{
		Handler: router,
		Addr:    addr,
	}
	return s
}

func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

func (s *Server) handleGoogleApp(c *gin.Context) {
	gid := c.Param("gid")
	log := s.Log.WithFields(logrus.Fields{
		"guildID": gid,
	})

	var payload AppInfo
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Errorf("Could not parse JSON body: %+v", err)
		return
	}

	if err := s.handleApplicant(gid, payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Errorf("Could handle applicant: %+v", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
	log.Infof("Handled applicant %+v", payload)
	return
}

func (s *Server) handleApplicant(gid string, app AppInfo) error {
	cName := app.ChannelName()

	cs, err := s.discord.GuildChannels(gid)
	if err != nil {
		return err
	}

	var (
		categ   *discordgo.Channel
		channel *discordgo.Channel
	)
	for i := range cs {
		ch := cs[i]
		s.Log.Debugf("Found channel [%d]: %s (type %s, parent %s)", i, ch.Name, ch.Type, ch.ParentID)
		switch {
		case ch.Type == discordgo.ChannelTypeGuildCategory && ch.Name == "applications": // TODO: don't hard-code this
			categ = ch
			if channel != nil && channel.ParentID != categ.ID {
				channel = nil
			}

		case ch.Type == discordgo.ChannelTypeGuildText && ch.Name == cName:
			if categ == nil || ch.ParentID == categ.ID {
				channel = ch
			}
		}
	}
	if categ == nil {
		return errors.Errorf("Did not find applications category")
	}
	if channel == nil {
		channel, err = s.discord.GuildChannelCreateComplex(gid, discordgo.GuildChannelCreateData{
			Name:     cName,
			Type:     discordgo.ChannelTypeGuildText,
			ParentID: categ.ID,
		})
		if err != nil {
			return errors.Wrapf(err, "could not create applicant channel '%s' for %+v", cName, app)
		}
		s.Log.Infof("Created channel %s in %s for %+v", cName, categ.Name, app)
	}
	_, err = s.discord.ChannelMessageSendEmbed(channel.ID, &discordgo.MessageEmbed{
		Title: app.Name,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Age", Value: app.Age, Inline: true},
			{Name: "Battle Tag", Value: app.BattleTag, Inline: true},
			{Name: "Armory", Value: app.ArmoryURL},
			{Name: "Logs", Value: app.LogsURL},
			{Name: "UI", Value: app.InterfaceURL},
		},
	})
	if err != nil {
		return errors.Wrapf(err, "could not send intro message for %+v to '%s'", app, cName)
	}
	for _, resp := range app.OtherResponses {
		parts := SplitWrap(resp.Answer, 1024)
		for i, part := range parts {
			embed := &discordgo.MessageEmbed{
				Title:       resp.Question,
				Description: part,
			}
			if len(parts) > 1 {
				embed.Footer = &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("(%d/%d)", i+1, len(parts)),
				}
			}
			_, err = s.discord.ChannelMessageSendEmbed(channel.ID, embed)
			if err != nil {
				return errors.Wrapf(err, "could not send reponse '%s' message %d for %+v to '%s'", resp.Question, i, app, cName)
			}
		}
	}
	return nil
}

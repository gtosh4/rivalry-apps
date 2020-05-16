package app

import (
	"strings"
	"time"
)

type (
	AppInfo struct {
		Timestamp time.Time `json:"timestamp"`

		Name         string `json:"name"`
		Age          string `json:"age"`
		BattleTag    string `json:"battle_tag"`
		ArmoryURL    string `json:"armory"`
		LogsURL      string `json:"logs"`
		InterfaceURL string `json:"ui"`

		OtherResponses []Response `json:"other"`
	}

	Response struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
	}
)

func (a *AppInfo) ChannelName() (n string) {
	if a.BattleTag != "" {
		n = a.BattleTag
		n = strings.ReplaceAll(n, "#", "-")
	} else {
		n = a.Name
		n = strings.ReplaceAll(n, " ", "-")
	}
	return n
}

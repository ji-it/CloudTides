package controller

import "github.com/robfig/cron/v3"

type VcdConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Org      string `json:"org"`
	Href     string `json:"href"`
	VDC      string `json:"vdc"`
	Token    string `json:"token"`
	Insecure bool   `json:"insecure"`
}

type Policy struct {
	CPU float64 `json:"cpu"`
	RAM float64 `json:"ram"`
}

var (
	letters  = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	cronjobs = map[uint]*cron.Cron{}
)

const (
	schedule      string = "*/5 * * * *"
	cleanSchedule string = "0 0 1 * *"
)

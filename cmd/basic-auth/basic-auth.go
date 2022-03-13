package main

import (
	"flag"

	"github.com/Kran001/basic-auth/internal/app"
)

const settingsFlagName = "settings"
const defaultSettingsLocation = "/etc/basic-auth/configs/settings.xml"

const sessionsSpoilTimeFlagName = "session-spoil-time"
const defaultSessionsSpoilTimeValue = "7|day"

func main() {
	serverSettings := flag.String(settingsFlagName, defaultSettingsLocation, "path to settings")
	sessionsValues := flag.String(sessionsSpoilTimeFlagName, defaultSessionsSpoilTimeValue, "sessions spoil values")

	flag.Parse()

	app.Run(*serverSettings, *sessionsValues)

}

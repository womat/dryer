package global

import (
	"dryer/pkg/dryer"
	"io"
	"time"
)

// VERSION holds the version information with the following logic in mind
//  1 ... fixed
//  0 ... year 2020, 1->year 2021, etc.
//  7 ... month of year (7=July)
//  the date format after the + is always the first of the month
//
// VERSION differs from semantic versioning as described in https://semver.org/
// but we keep the correct syntax.
//TODO: increase version number to 1.0.1+2020xxyy
const VERSION = "0.0.1+20201227"
const MODULE = "dryer"

type DebugConf struct {
	File io.WriteCloser
	Flag int
}

type WebserverConf struct {
	Port        int
	Webservices map[string]bool
}

type Configuration struct {
	DataCollectionInterval time.Duration
	BackupInterval         time.Duration
	DataFile               string
	Debug                  DebugConf
	Webserver              WebserverConf
	MeterURL               string
}

// Config holds the global configuration
var Config Configuration

// Measurements hold all measured  heat pump values
var Measurements *dryer.Measurements

func init() {
	Config = Configuration{Webserver: WebserverConf{Webservices: map[string]bool{}}}
}

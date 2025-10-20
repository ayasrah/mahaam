package conf

import (
	"encoding/json"
	"log"
	"os"
)

func NewConf(fileName string) *Conf {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error loading config file" + err.Error())
	}
	decoder := json.NewDecoder(file)
	cfg := Conf{}
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal("Error loading config file" + err.Error())
	}
	return &cfg
}

type Conf struct {
	ApiName                     string
	ApiVersion                  string
	EnvName                     string
	DBUrl                       string
	LogFile                     string
	LogFileSizeLimit            int
	LogFileCountLimit           int
	LogFileOutputTemplate       string
	LogFileRollingInterval      string
	HTTPPort                    int
	TokenSecretKey              string
	EmailAccountSID             string
	EmailVerificationServiceSID string
	EmailAuthToken              string
	TestEmails                  []string
	TestSID                     string
	TestOTP                     string
	LogReqEnabled               bool
}

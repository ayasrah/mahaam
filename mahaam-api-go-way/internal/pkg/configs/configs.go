package configs

import (
	"encoding/json"
	"log"
	"os"
)

var (
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
)

func InitConfigs() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	conf := config{}
	err := decoder.Decode(&conf)
	if err != nil {
		log.Fatal("Error loading config file" + err.Error())
	}
	ApiName = conf.ApiName
	ApiVersion = conf.ApiVersion
	EnvName = conf.EnvName
	DBUrl = conf.DBUrl
	LogFile = conf.LogFile
	LogFileSizeLimit = conf.LogFileSizeLimit
	LogFileCountLimit = conf.LogFileCountLimit
	LogFileOutputTemplate = conf.LogFileOutputTemplate
	LogFileRollingInterval = conf.LogFileRollingInterval
	HTTPPort = conf.HTTPPort
	TokenSecretKey = conf.TokenSecretKey
	EmailAccountSID = conf.EmailAccountSID
	EmailVerificationServiceSID = conf.EmailVerificationServiceSID
	EmailAuthToken = conf.EmailAuthToken
	TestEmails = conf.TestEmails
	TestSID = conf.TestSID
	TestOTP = conf.TestOTP
	LogReqEnabled = conf.LogReqEnabled
}

type config struct {
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

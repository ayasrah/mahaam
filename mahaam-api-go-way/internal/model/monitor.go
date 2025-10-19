package model

import "github.com/google/uuid"

type Traffic struct {
	ID       uuid.UUID `json:"id" db:"id"`
	HealthID uuid.UUID `json:"healthId" db:"health_id"`
	Method   string    `json:"method" db:"method"`
	Path     string    `json:"path" db:"path"`
	Code     int       `json:"code" db:"code"`
	Elapsed  int64     `json:"elapsed" db:"elapsed"`
	Headers  string    `json:"headers" db:"headers"`
	Request  string    `json:"request" db:"request"`
	Response string    `json:"response" db:"response"`
}

type TrafficHeaders struct {
	UserID     uuid.UUID `json:"userId"`
	DeviceID   uuid.UUID `json:"deviceId"`
	AppVersion string    `json:"appVersion"`
	AppStore   string    `json:"appStore"`
}

type Health struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ApiName    string    `json:"apiName" db:"api_name"`
	ApiVersion string    `json:"apiVersion" db:"api_version"`
	NodeIP     string    `json:"nodeIp" db:"node_ip"`
	NodeName   string    `json:"nodeName" db:"node_name"`
	EnvName    string    `json:"envName" db:"env_name"`
}

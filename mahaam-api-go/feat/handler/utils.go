package handler

import (
	"github.com/google/uuid"
)

func ExtractMeta(c Ctx) Meta {
	return Meta{
		UserID:   ExtractUserID(c),
		DeviceID: ExtractDeviceID(c),
	}
}

func ExtractUserID(c Ctx) UUID {
	userId, ok := c.Value("userId").(UUID)
	if !ok || userId == uuid.Nil {
		panic("userId not found in context")
	}
	return userId
}

func ExtractDeviceID(c Ctx) UUID {
	deviceId, ok := c.Value("deviceId").(UUID)
	if !ok || deviceId == uuid.Nil {
		panic("deviceId not found in context")
	}
	return deviceId
}

func ExtractTrafficID(c Ctx) UUID {
	trafficID, ok := c.Value("trafficID").(UUID)
	if !ok || trafficID == uuid.Nil {
		panic("trafficID not found in context")
	}
	return trafficID
}

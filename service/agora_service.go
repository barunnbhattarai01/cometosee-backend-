package service

import (
	"time"

	rtctokenbuilder "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
)

type AgoraService interface {
	GenerateToken(channel string, uid uint32, expire int) (string, error)
}

type agoraService struct {
	appID   string
	appCert string
}

func NewAgoraService(appID, appCert string) AgoraService {
	return &agoraService{
		appID:   appID,
		appCert: appCert,
	}
}

func (a *agoraService) GenerateToken(channel string, uid uint32, expire int) (string, error) {

	expireTime := uint32(time.Now().Unix() + int64(expire))

	token, err := rtctokenbuilder.BuildTokenWithUid(
		a.appID,
		a.appCert,
		channel,
		uid,
		rtctokenbuilder.RolePublisher,
		expireTime,
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

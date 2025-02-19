package model

import (
	"time"
)

type DIDDocument struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	PublicKey      []string  `json:"publicKey" gorm:"type:json"`
	Authentication []string  `json:"authentication" gorm:"type:json"`
	Created        time.Time `json:"created"`
	Updated        time.Time `json:"updated"`
}

func (DIDDocument) TableName() string {
	return "did_documents"
}

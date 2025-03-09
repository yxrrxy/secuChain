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

//type DID struct {
//	ID       string      `json:"id" gorm:"primaryKey;type:varchar(255)"`
//	AuthKey  KeyPairs    `json:"authKey" gorm:"type:jsonb"`
//	RecyKey  KeyPairs    `json:"recyKey" gorm:"type:jsonb"`
//	Document DIDDocument `json:"document" gorm:"foreignKey:ID"`
//}

//func (DID) TableName() string {
//	return "did_documents"
//}
//

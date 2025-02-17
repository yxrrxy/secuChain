package model

// SPDXSBOM represents SPDX SBOM data structure
type SPDXSBOM struct {
	SPDXID       string       `json:"spdxID"`
	Name         string       `json:"name"`
	VersionInfo  string       `json:"versionInfo"`
	Supplier     string       `json:"supplier"`
	ExternalRefs ExternalRefs `json:"externalRefs"` // 修改为使用独立的结构体类型
}

// ExternalRefs 是一个独立的结构体类型，用于表示 SPDX SBOM 中的 externalRefs 字段
type ExternalRefs struct {
	ReferenceCategory string `json:"referenceCategory"`
	ReferenceLocator  string `json:"referenceLocator"`
	ReferenceType     string `json:"referenceType"`
}

// CDXSBOM represents CycloneDX SBOM data structure
type CDXSBOM struct {
	BomRef  string `json:"bom-ref"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Purl    string `json:"purl"`
}

// SBOM represents a generic SBOM structure that can encapsulate either SPDXSBOM or CDXSBOM
type SBOM struct {
	ID       string      `json:"id" gorm:"primaryKey"`
	DID      string      `json:"did" gorm:"index"` // Assuming DID is a required field
	SPDXSBOM *SPDXSBOM   `json:"spdxSBOM,omitempty"`
	CDXSBOM  *CDXSBOM    `json:"cdxSBOM,omitempty"`
}

func (SBOM) TableName() string {
	return "sboms"
}
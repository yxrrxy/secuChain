package model

type SPDXSBOM struct {
	SPDXID       string       `json:"spdxID"`
	Name         string       `json:"name"`
	VersionInfo  string       `json:"versionInfo"`
	Supplier     string       `json:"supplier"`
	ExternalRefs externalRefs `json:"externalRefs"`
}

type externalRefs struct {
	ReferenceCategory string `json:"referenceCategory"`
	ReferenceLocator  string `json:"referenceLocator"`
	ReferenceType     string `json:"referenceType"`
}

type CDXSBOM struct {
	BomRef  string `json:"bom-ref"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Purl    string `json:"purl"`
}

func (SPDXSBOM) TableName() string {
	return "spdxsboms"
}

func (CDXSBOM) TableName() string {
	return "cdxsboms"
}

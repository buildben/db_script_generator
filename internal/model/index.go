package model

type IndexFile struct {
	Version string   `json:"version"`
	Init    []string `json:"init"`
	Schemas []Schema `json:"schemas"`
}

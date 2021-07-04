package model

type IndexFile struct {
	Version int64    `json:"version"`
	Init    []string `json:"init"`
	Schemas []Schema `json:"schemas"`
}

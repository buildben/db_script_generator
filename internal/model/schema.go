package model

type Schema struct {
	Name     string   `json:"name"`
	Types    []string `json:"types"`
	Tables   []string `json:"tables"`
	Data     []string `json:"data"`
	Views    []string `json:"views"`
	Routines []string `json:"routines"`
}

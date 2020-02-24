package models

type ValidField struct {
	Field string `json:"field"  example:"id"`
	Tag   string `json:"tag" example:"min"`
	Param string `json:"param" example:"1"`
}

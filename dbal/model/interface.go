package model

// Setter the model setter
type Setter interface {
	GetAttr(name string) *Attribute
	SetAttr(name string, attr Attribute)
	Fill(attributes interface{}, v ...interface{}) *Model
	Set(name string, value interface{}, v ...interface{}) *Model
}

package model

// Setter the model setter
type Setter interface {
	GetAttr(name string) *Attribute
	SetAttr(name string, attr Attribute)
}

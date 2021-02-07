package grammar

// Grammar the Grammar inteface
type Grammar interface {
	Exists(table Table) string
}

// Table the table struct
type Table struct {
	Name    string
	Comment string
}

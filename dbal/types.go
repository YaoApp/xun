package dbal

// Schema The database Schema interface
type Schema interface {
	Create()
	Drop()
	DropIfExists()
	Rename()
	Primary()

	BigInteger()
	String()
}

// Blueprint the dbal driver
type Blueprint struct{ Schema }

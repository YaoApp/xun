package schema

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

// Blueprint the dbal schema driver
type Blueprint struct{ Schema }

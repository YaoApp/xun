package schema

// Schema The schema interface
type Schema interface {
	Get(string) (Blueprint, error)
	Create(string, func(table Blueprint)) error
	Drop(string) error
	Alter(string, func(table Blueprint)) error
	HasTable(string) bool
	Rename(string, string) error
	DropIfExists(string) error

	MustGet(string) Blueprint
	MustCreate(string, func(table Blueprint)) Blueprint
	MustDrop(string)
	MustAlter(string, func(table Blueprint)) Blueprint
	MustRename(string, string) Blueprint
	MustDropIfExists(string)
}

// Blueprint the table operating interface
type Blueprint interface {
	GetName() string
	GetColumns() map[string]*Column
	GetIndexes() map[string]*Index

	GetColumn(name string) *Column
	HasColumn(name ...string) bool
	RenameColumn(old string, new string) *Column
	DropColumn(name ...string)

	GetPrimary() *Primary
	CreatePrimary(columnName ...string)
	CreatePrimaryWithName(name string, columnName ...string)
	DropPrimary()

	GetIndex(name string) *Index
	HasIndex(name ...string) bool
	CreateIndex(key string, columnNames ...string)
	CreateUnique(key string, columnNames ...string)
	RenameIndex(old string, new string) *Index
	DropIndex(key ...string)

	CreateUniqueConstraint(name string, columnNames ...string)
	GetUniqueConstraint(name string)
	DropUniqueConstraint(name string)

	String(name string, length int) *Column
	BigInteger(name string) *Column
	UnsignedBigInteger(name string) *Column
	BigIncrements(name string) *Column
	ID(name string) *Column
}

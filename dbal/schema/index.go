package schema

// Drop mark as dropped for the index
func (index *Index) Drop() {
	index.Dropped = true
}

// Rename mark as renamed with the given name for the index
func (index *Index) Rename(new string) {
	index.Newname = new
}

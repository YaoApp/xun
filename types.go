package xun

// R alias map[string]interface{}, R is the first letter of "Row"
type R map[string]interface{}

// N an numberic value,  R is the first letter of "Numberic"
type N struct {
	Number interface{}
}

// T an datetime value, T is the first letter of "Time"
type T struct {
	Time interface{}
}

// P an Paginator struct, P is the first letter of "Paginator"
type P struct {
	Items        []interface{}
	Total        int
	PageCount    int
	Perpage      int
	CurrentPage  int
	NextPage     int
	PreviousPage int
	LastPage     int
	Options      map[string]interface{}
}

package xun

// R alias map[string]interface{}, R is the first letter of "Row"
type R map[string]interface{}

// N an numberic value,  R is the first letter of "Numberic"
type N struct {
	Value interface{}
}

// T an datetime value, T is the first letter of "Time"
type T struct {
	Value interface{}
}

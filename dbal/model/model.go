package model

// Search search by given params
func (model *Model) Search() interface{} {
	return nil
}

// Find find by primary key
func (model *Model) Find(v ...interface{}) interface{} {
	if len(v) > 0 {
		return v[0]
	}
	return model
}

// Flow process a flow by the given flow name and return the result
func (model *Model) Flow(name string) interface{} {
	return nil
}

// FlowRaw process a flow by the given json description file and return the result
func (model *Model) FlowRaw(flow []byte) interface{} {
	return nil
}

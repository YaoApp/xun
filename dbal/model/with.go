package model

import (
	"fmt"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/utils"
)

// With where the array key is a relationship name and the array value is a closure that adds additional constraints to the eager loading query
func (model *Model) With(args ...interface{}) *Model {
	name, closure := prepareWithArgs(args...)
	var rel *Relationship = nil
	if attr, has := model.attributes[name]; has {
		rel = attr.Relationship
	}

	if rel == nil {
		invalidRelationship()
	}

	switch rel.Type {
	case "hasOne":
		model.withHas("hasOne", rel, name, closure)
		break
	case "hasMany":
		model.withHas("hasMany", rel, name, closure)
		break
	case "hasOneThrough":
		model.withThrough("hasOneThrough", rel, name, closure)
	}
	return model
}

// withHasThrough
func (model *Model) withThrough(typ string, rel *Relationship, name string, closure func(query.Query)) {

	length := len(rel.Models)
	if length < 2 {
		invalidRelationship()
	}

	linkLength := len(rel.Links)
	if linkLength != length*2 {
		invalidRelationship(fmt.Sprintf("Links should have %d fields", length*2))
	}

	model.
		MakeModelForRelationship(rel.Models[length-1]).
		BasicQueryForRelationship(rel.Columns, closure)

	return
}

// withHasMany
func (model *Model) withHas(typ string, rel *Relationship, name string, closure func(query.Query)) {

	if len(rel.Models) < 1 {
		invalidRelationship()
	}

	relModel := model.
		MakeModelForRelationship(rel.Models[0]).
		BasicQueryForRelationship(rel.Columns, closure)

	// links M1.Local, (->) M2.Foreign, M2.Local, (->) M3.Foreign ...
	local, foreign := relModel.GetRelationshipLink(rel.Models[0], rel.Links...)
	relModel.SelectAddColumn(foreign)

	model.withs = append(model.withs, With{
		Name:    name,
		Type:    typ,
		Query:   relModel.QueryBuilder(),
		Local:   local,
		Foreign: foreign,
	})

}

// ExecuteWiths Execute the withs query and merge result
func (model *Model) ExecuteWiths(rows []xun.R) error {

	if len(rows) == 0 {
		return nil
	}

	for _, with := range model.withs {
		var err error
		switch with.Type {
		case "hasOne":
			err = model.executeHasOne(rows, with)
		case "hasMany":
			err = model.executeHasMany(rows, with)
		}
		if err != nil {
			rows = nil
			return err
		}
	}
	return nil
}

func (model *Model) executeHasMany(rows []xun.R, with With) error {

	ids := []interface{}{}
	for _, row := range rows {
		id := row.Get(with.Local)
		if id == nil {
			return fmt.Errorf("The %s is nil", with.Local)
		}
		ids = append(ids, id)
	}

	// get the relation rows
	relrows, err := with.Query.WhereIn(with.Foreign, ids).Get()
	if err != nil {
		return err
	}

	// merge the result
	relRowsMap := map[interface{}][]xun.R{}
	for _, rel := range relrows {
		id := rel.Get(with.Foreign)
		if _, has := relRowsMap[id]; !has {
			relRowsMap[id] = []xun.R{}
		}
		relRowsMap[id] = append(relRowsMap[id], rel)
	}

	for _, row := range rows {
		id := row.Get(with.Local)
		if rel, has := relRowsMap[id]; has {
			row[with.Name] = rel
		} else {
			row[with.Name] = []xun.R{}
		}
	}

	return nil
}

func (model *Model) executeHasOne(rows []xun.R, with With) error {

	ids := []interface{}{}
	for _, row := range rows {
		id := row.Get(with.Local)
		if id == nil {
			return fmt.Errorf("The %s is nil", with.Local)
		}
		ids = append(ids, id)
		ids = utils.InterfaceUnique(ids)
	}

	// get the relation rows
	relrows, err := with.Query.WhereIn(with.Foreign, ids).Get()
	if err != nil {
		return err
	}

	// merge the result
	relRowsMap := map[interface{}]xun.R{}
	for _, rel := range relrows {
		id := rel.Get(with.Foreign)
		relRowsMap[id] = rel

	}
	for _, row := range rows {
		id := row.Get(with.Local)
		if rel, has := relRowsMap[id]; has {
			row[with.Name] = rel
		} else {
			row[with.Name] = xun.MakeRow()
		}
	}

	return nil
}

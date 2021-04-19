package model

import (
	"fmt"
	"strings"

	"github.com/yaoapp/xun"
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
	case "hasManyThrough":
		model.withThrough("hasManyThrough", rel, name, closure)
	}
	return model
}

// withHas hasOne/hasMany
func (model *Model) withHas(typ string, rel *Relationship, name string, closure func(Basic)) {

	if len(rel.Models) < 1 {
		invalidRelationship()
	}

	relModel := model.
		MakeModelForRelationship(rel.Models[0]).
		BasicQueryForRelationship(rel.Columns, closure)

	links := model.parseLinks(rel.Models, rel.Links)
	relModel.SelectAddColumn(fmt.Sprintf("%s as temp_foreign_id", links[0].ForeignKey), relModel.GetTableName())
	relModel.QueryBuilder().Distinct("temp_foreign_id")
	model.withs = append(model.withs, With{
		Name:  name,
		Type:  typ,
		Basic: relModel,
		Links: links,
	})

}

// withHasThrough hasOneThrough/hasManyThrough
func (model *Model) withThrough(typ string, rel *Relationship, name string, closure func(Basic)) {

	links := model.parseLinks(rel.Models, rel.Links)
	length := len(links)
	if length < 2 {
		invalidRelationship("the hasThrough type relationship should have two pair links at least.")
	}

	relModel := model.
		MakeModelForRelationship(rel.Models[0]).
		BasicQueryForRelationship(rel.Columns, closure)

	// Joins
	for i := 1; i < len(links); i++ {
		link := links[i]
		relModel.
			QueryBuilder().
			LeftJoin(fmt.Sprintf("%s as %s", link.To, link.AliasTo), link.Local, "=", link.Foreign)
	}

	// the last link
	link := links[length-1]
	relModel.SelectAddColumn(fmt.Sprintf("%s as temp_foreign_id", links[0].Foreign))
	relModel.QueryBuilder().WhereNotNull(link.Foreign).Distinct("temp_foreign_id")
	relModel.TableColumnize(link.AliasTo)

	model.withs = append(model.withs, With{
		Name:  name,
		Type:  typ,
		Basic: relModel,
		Links: links,
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
		case "hasOneThrough":
			err = model.executeHasOneThrough(rows, with)
		case "hasManyThrough":
			err = model.executeHasManyThrough(rows, with)
		}
		if err != nil {
			rows = nil
			return err
		}
	}
	return nil
}

func (model *Model) executeHasOne(rows []xun.R, with With) error {

	foreign := with.Links[0].ForeignKey
	local := with.Links[0].LocalKey

	ids := []interface{}{}
	for _, row := range rows {
		id := row.Get(local)
		if id == nil {
			return fmt.Errorf("The %s is nil", local)
		}
		ids = append(ids, id)
		ids = utils.InterfaceUnique(ids)
	}

	// get the relation rows
	relationRows, err := with.Basic.QueryBuilder().WhereIn(foreign, ids).Get()
	if err != nil {
		return err
	}

	// merge the result
	model.mergeOne(with.Name, local, rows, relationRows)

	return nil
}

func (model *Model) executeHasMany(rows []xun.R, with With) error {

	foreign := with.Links[0].ForeignKey
	local := with.Links[0].LocalKey
	ids := []interface{}{}
	for _, row := range rows {
		id := row.Get(local)
		if id == nil {
			return fmt.Errorf("The %s is nil", local)
		}
		ids = append(ids, id)
	}

	// get the relation rows
	relationRows, err := with.Basic.
		QueryBuilder().
		WhereIn(foreign, ids).
		Get()
	if err != nil {
		return err
	}

	model.mergeMany(with.Name, local, rows, relationRows)
	return nil
}

func (model *Model) executeHasOneThrough(rows []xun.R, with With) error {

	foreign := with.Links[0].Foreign
	local := with.Links[0].LocalKey

	ids := []interface{}{}
	for _, row := range rows {
		id := row.Get(local)
		if id == nil {
			return fmt.Errorf("The %s is nil", local)
		}
		ids = append(ids, id)
		ids = utils.InterfaceUnique(ids)
	}

	// get the relation rows
	relationRows, err := with.Basic.QueryBuilder().WhereIn(foreign, ids).Get()
	if err != nil {
		return err
	}

	model.mergeOne(with.Name, local, rows, relationRows)
	return nil
}

func (model *Model) executeHasManyThrough(rows []xun.R, with With) error {

	foreign := with.Links[0].Foreign
	local := with.Links[0].LocalKey

	ids := []interface{}{}
	for _, row := range rows {
		id := row.Get(local)
		if id == nil {
			return fmt.Errorf("The %s is nil", local)
		}
		ids = append(ids, id)
		ids = utils.InterfaceUnique(ids)
	}

	// get the relation rows
	relationRows, err := with.Basic.
		QueryBuilder().
		WhereIn(foreign, ids).
		Get()

	if err != nil {
		return err
	}

	model.mergeMany(with.Name, local, rows, relationRows)
	return nil
}

func (model *Model) mergeOne(name string, local string, rows []xun.R, relationRows []xun.R) {

	relRowsMap := map[interface{}]xun.R{}
	for _, rel := range relationRows {
		id := rel.Get("temp_foreign_id")
		rel.Del("temp_foreign_id")
		relRowsMap[id] = rel
	}

	for _, row := range rows {
		id := row.Get(local)
		if rel, has := relRowsMap[id]; has {
			row[name] = rel
		} else {
			row[name] = xun.MakeRow()
		}
	}

}

func (model *Model) mergeMany(name string, local string, rows []xun.R, relationRows []xun.R) {
	// merge the result
	relRowsMap := map[interface{}][]xun.R{}
	for _, rel := range relationRows {
		id := rel.Get("temp_foreign_id")
		rel.Del("temp_foreign_id")
		if _, has := relRowsMap[id]; !has {
			relRowsMap[id] = []xun.R{}
		}
		relRowsMap[id] = append(relRowsMap[id], rel)
	}

	for _, row := range rows {
		id := row.Get(local)
		if rel, has := relRowsMap[id]; has {
			row[name] = rel
		} else {
			row[name] = []xun.R{}
		}
	}
}

// parseLinks make the links with given relations
func (model *Model) parseLinks(models []string, keys []string) []Link {

	length := len(models)
	if length < 1 {
		invalidRelationship()
	}

	keyLength := len(keys)
	if keyLength != length*2 {
		invalidRelationship(fmt.Sprintf("Links should have %d fields", length*2))
	}

	aliasFrom := model.GetName()
	from := model.GetTableName()
	links := []Link{}
	for i, name := range models {
		rel := model.MakeModelForRelationship(name)
		to := rel.GetTableName()
		aliasTo := rel.GetName()
		localKey := keys[i*2]
		foreignKey := keys[i*2+1]
		local := localKey
		foreign := foreignKey
		if !strings.Contains(localKey, ".") {
			local = fmt.Sprintf("%s.%s", aliasFrom, localKey)
		}
		if !strings.Contains(foreignKey, ".") {
			foreign = fmt.Sprintf("%s.%s", aliasTo, foreignKey)
		}
		link := Link{
			From:       from,
			To:         to,
			AliasFrom:  aliasFrom,
			AliasTo:    aliasTo,
			LocalKey:   localKey,
			ForeignKey: foreignKey,
			Local:      local,
			Foreign:    foreign,
		}
		links = append(links, link)
		from = to
		aliasFrom = aliasTo
	}
	return links
}

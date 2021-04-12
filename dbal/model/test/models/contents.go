package models

// SchemaFileContents the json shcema files
var SchemaFileContents = map[string][]byte{

	"models/member.json": []byte(`{
		"name": "Member",
		"table": {
		  "name": "member",
		  "comment": "Member",
		  "engine": "InnoDB"
		},
		"columns": [
		  {
			"comment": "Member ID",
			"name": "id",
			"type": "ID",
			"example": 10028,
			"title": "Member ID",
			"description": "The member ID",
			"validation": {
			  "pattern": "^[0-9]{1,16}$",
			  "description": "Member ID must be the integer"
			}
		  }
		]
	  }
	`),

	"models/user.json": []byte(`{
		"name": "User",
		"table": {
		  "name": "user",
		  "comment": "User",
		  "engine": "InnoDB"
		}
	  }
	`),
}

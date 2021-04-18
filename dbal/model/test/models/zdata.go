// Package models  Sun Apr 18 16:19:40 CST 2021
// THIS FILE IS AUTO-GENERATED DO NOT MODIFY MANUALLY
package models

// SchemaFileContents the json shcema files
var SchemaFileContents = map[string][]byte{

	"models/car.json": []byte(`
{
  "name": "Car",
  "table": {
    "name": "car",
    "comment": "The cars",
    "engine": "InnoDB"
  },
  "columns": [
    { "name": "id", "type": "ID", "comment": "Car ID" },
    { "name": "name", "type": "string", "length": 200, "comment": "Car Name" },
    { "name": "manu_id", "type": "bigInteger", "comment": "Manufacturer id" }
  ],
  "relationships": [
    {
      "name": "manu",
      "type": "hasOne",
      "models": ["manu"],
      "links": ["manu_id", "id"],
      "columns": ["name", "type", "intro"]
    },
    {
      "name": "users",
      "type": "hasManyThrough",
      "models": ["user_car", "user"],
      "links": ["id", "car_id", "user_id", "id"]
    }
  ],
  "indexes": [{ "name": "manu_id", "columns": ["manu_id"] }],
  "option": {
    "soft_deletes": true,
    "timestamps": true
  },
  "values": [
    { "name": "Tesla Model 3", "manu_id": 1, "deleted_at": null },
    {
      "name": "Tesla Cybertruck",
      "manu_id": 1,
      "deleted_at": "2020-03-04 22:24:39"
    },
    { "name": "Tesla Model X", "manu_id": 1, "deleted_at": null },
    { "name": "SERENA E-POWER", "manu_id": 3, "deleted_at": null },
    { "name": "Teana", "manu_id": 3, "deleted_at": null },
    { "name": "RAV4", "manu_id": 2, "deleted_at": null },
    { "name": "CH-R", "manu_id": 2, "deleted_at": null }
  ]
}
`),
	"models/manu.json": []byte(`
{
  "name": "Manu",
  "table": {
    "name": "manu",
    "comment": "The Manufacturer",
    "engine": "InnoDB"
  },
  "columns": [
    { "name": "id", "type": "ID", "comment": "Manufacturer ID" },
    {
      "name": "name",
      "type": "string",
      "length": 200,
      "comment": "Manufacturer Name"
    },
    {
      "name": "intro",
      "type": "text",
      "comment": "Manufacturer introduction"
    },
    {
      "name": "type",
      "type": "enum",
      "option": ["electric", "gas", "mix"],
      "comment": "Energy type",
      "index": true
    }
  ],
  "relationships": [
    {
      "name": "cars",
      "type": "hasMany",
      "models": ["car"],
      "links": ["id", "manu_id"]
    },
    {
      "name": "user",
      "type": "hasOneThrough",
      "models": ["car", "user_car", "user"],
      "links": ["id", "manu_id", "car_id", "car_id", "user_id", "id"]
    },
    {
      "name": "users",
      "type": "hasManyThrough",
      "models": ["car", "user_car", "user"],
      "links": ["id", "manu_id", "car_id", "car_id", "user_id", "id"]
    }
  ],
  "values": [
    {
      "name": "Tesla",
      "type": "electric",
      "intro": "Tesla is accelerating the world's transition to sustainable energy with electric cars, solar and integrated renewable energy solutions for homes and businesses."
    },
    {
      "name": "Toyota",
      "type": "gas",
      "intro": "The Toyota Motor Corporation is a Japanese multinational automotive manufacturer headquartered in Toyota, Aichi, Japan. "
    },
    {
      "name": "Nissan",
      "type": "mix",
      "intro": "The Nissan Motor Company, Ltd. is a Japanese multinational automobile manufacturer headquartered in Nishi-ku, Yokohama, Japan."
    }
  ]
}
`),
	"models/member.flow.json": []byte(`
{
  "name": "vip",
  "method": {
    "name": "getVip",
    "in": [{ "name": "id", "required": true }],
    "out": [{ "name": "field", "mapping": "field::name" }]
  }
}
`),
	"models/member.json": []byte(`
{
  "name": "Member",
  "table": {
    "name": "member",
    "comment": "Member",
    "engine": "InnoDB"
  },
  "columns": [
    {
      "name": "id",
      "comment": "Member ID",
      "type": "ID",
      "title": "Member ID",
      "description": "The member ID",
      "validation": {
        "pattern": "^[0-9]{1,16}$",
        "description": "Member ID must be the integer"
      },
      "example": 10001
    },
    {
      "name": "user_id",
      "comment": "User ID, 1v1",
      "type": "bigInteger",
      "nullable": false,
      "unique": true,
      "example": 20001
    },
    {
      "name": "name",
      "comment": "Real Name",
      "type": "string",
      "length": 80,
      "example": "John"
    },
    {
      "name": "score",
      "comment": "The member score",
      "type": "float",
      "nullable": false,
      "default": 0.0,
      "precision": 5,
      "scale": 2,
      "example": 3.28
    },
    {
      "name": "level",
      "comment": "The member level",
      "type": "enum",
      "option": ["silver", "gold"],
      "default": "silver",
      "example": "gold"
    },
    {
      "name": "expired_at",
      "comment": "The member expired at",
      "type": "timestamp",
      "nullable": false,
      "default_raw": "NOW()",
      "example": "2021-04-13 22:19:62"
    }
  ],
  "relationships": [
    {
      "name": "user",
      "type": "hasOne",
      "models": ["User"],
      "links": ["user_id", "id"]
    },
    {
      "name": "cars",
      "type": "hasManyThrough",
      "models": ["UserCar", "Car"],
      "links": ["user_id", "user_id", "car_id", "id"]
    }
  ]
}
`),
	"models/null.json": []byte(`
{}
`),
	"models/user.flow.json": []byte(`
{
  "name": "login",
  "method": {
    "name": "login",
    "in": [
      { "name": "email", "field": "email", "required": true },
      { "name": "password", "type": "string", "required": true }
    ],
    "out": [
      { "name": "name", "field": "name" },
      { "name": "mobile", "field": "mobile" }
    ]
  }
}
`),
	"models/user.json": []byte(`
{
  "name": "User",
  "table": {
    "name": "user",
    "comment": "User",
    "engine": "InnoDB"
  },
  "columns": [
    { "name": "id" },
    { "name": "nickname", "comment": "The user nick name", "unique": true },
    { "name": "bio", "type": "text", "comment": "The user bio" },
    { "name": "gender", "type": "tinyInteger", "default": 3, "index": true },
    { "name": "address", "length": 300, "comment": "The address" },
    { "name": "score", "precision": 5, "scale": 2, "index": true },
    { "name": "status", "default": "WAITING", "index": true }
  ],
  "relationships": [
    {
      "name": "member",
      "type": "hasOne",
      "models": ["Member"],
      "links": ["id", "user_id"]
    },
    {
      "name": "cars",
      "type": "hasManyThrough",
      "models": ["UserCar", "Car"],
      "links": ["id", "user_id", "car_id", "id"]
    }
  ],
  "values": [
    {
      "nickname": "admin",
      "bio": "the default adminstor",
      "gender": 1,
      "vote": 0,
      "score": 1.25,
      "address": "default path",
      "status": "DONE"
    }
  ]
}
`),
	"models/user_car.json": []byte(`
{
  "name": "UserCar",
  "table": {
    "name": "user_car",
    "comment": "User cars",
    "engine": "InnoDB"
  },
  "columns": [
    { "name": "id", "type": "ID", "comment": "Use Car ID" },
    { "name": "car_id", "type": "bigInteger", "index": true },
    { "name": "user_id", "type": "bigInteger", "index": true }
  ],
  "indexes": [{ "name": "car_id_user_id", "columns": ["car_id", "user_id"] }]
}
`),
}

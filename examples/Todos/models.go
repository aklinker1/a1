package main

import (
	"fmt"

	a1 "github.com/aklinker1/a1/pkg"
	a1Types "github.com/aklinker1/a1/pkg/types"
)

func postgreSQLTable(tableName string) a1.DataLoaderConfig {
	return a1.DataLoaderConfig{
		Source: "PostgreSQL",
		Group:  tableName,
	}
}

var models = a1.ModelMap{
	"User": a1.Model{
		DataLoader: postgreSQLTable("users"),
		GraphQL: a1.GraphQLConfig{
			DisableCreate: true,
			DisableUpdate: true,
			DisableDelete: true,
		},
		Fields: a1.FieldMap{
			"_id": a1.Field{
				Type:       a1Types.ID,
				PrimaryKey: true,
			},
			"username": a1Types.String,
			"email": a1.Field{
				Type:   "Email",
				Hidden: true,
			},
			"passwordHash": a1.Field{
				Type:      a1Types.String,
				Hidden:    true,
				DataField: "password_hash",
			},
		},
	},
	"MyUser": a1.Model{
		Extends: "User",
		GraphQL: a1.GraphQLConfig{
			DisableSelectMultiple: true,
			DisableCreate:         true,
			DisableUpdate:         true,
			DisableDelete:         true,
		},
		Fields: a1.FieldMap{
			"email":      "Email",
			"validation": "Validation",
			"firstName":  a1Types.String,
			"lastName":   a1Types.String,

			"fullName": a1.VirtualField{
				Type:           a1Types.String,
				RequiredFields: []string{"firstName", "lastName"},
				Compute: func(data map[string]interface{}) interface{} {
					return fmt.Sprintf("%s %s", data["firstName"], data["lastName"])
				},
			},

			"todos": a1.LinkedField{
				LinkedModel: "Todo",
				Type:        a1.OneToMany,
				Field:       "_id",
				LinkedField: "_userId",
			},
			"preferences": a1.LinkedField{
				LinkedModel: "Preferences",
				Type:        a1.OneToOne,
				Field:       "_id",
				LinkedField: "_userId",
			},
		},
	},
	"Preferences": a1.Model{
		DataLoader: postgreSQLTable("preferences"),
		GraphQL: a1.GraphQLConfig{
			DisableSelectMultiple: true,
			DisableCreate:         true,
			DisableDelete:         true,
		},
		Fields: a1.FieldMap{
			"_id": a1.Field{
				Type:       a1Types.ID,
				PrimaryKey: true,
			},
			"_userId": a1.Field{
				Type:      a1Types.ID,
				DataField: "user_id",
			},
			"theme": "Theme",

			"user": a1.LinkedField{
				LinkedModel: "MyUser",
				Type:        a1.OneToOne,
				Field:       "_userId",
				LinkedField: "_id",
			},
		},
	},
	"Todo": a1.Model{
		DataLoader: postgreSQLTable("todos"),
		Fields: a1.FieldMap{
			"_id": a1.Field{
				Type:       a1Types.ID,
				PrimaryKey: true,
			},
			"_userId": a1.Field{
				Type:      a1Types.ID,
				DataField: "user_id",
			},
			"message": a1Types.String,
			"isCompleted": a1.Field{
				Type:      a1Types.Bool,
				DataField: "is_completed",
			},

			"user": a1.LinkedField{
				LinkedModel: "User",
				Type:        a1.ManyToOne,
				Field:       "_userId",
				LinkedField: "_id",
			},
			"tags": a1.LinkedField{
				LinkedModel: "TodoTag",
				Type:        a1.OneToMany,
				Field:       "_id",
				LinkedField: "_todoId",
			},
		},
	},
	"Tag": a1.Model{
		DataLoader: postgreSQLTable("tags"),
		Fields: a1.FieldMap{
			"_name": a1.Field{
				Type:       a1Types.String,
				PrimaryKey: true,
			},
			"addedAt": a1Types.Date,

			"todoTags": a1.LinkedField{
				LinkedModel: "TodoTag",
				Type:        a1.OneToMany,
				Field:       "_name",
				LinkedField: "_tagName",
			},
		},
	},
	"TodoTag": a1.Model{
		DataLoader: postgreSQLTable("todo_tags"),
		GraphQL: a1.GraphQLConfig{
			DisableSelectOne:      true,
			DisableSelectMultiple: true,
			DisableCreate:         true,
			DisableUpdate:         true,
			DisableDelete:         true,
		},
		Fields: a1.FieldMap{
			"_id": a1.Field{
				Type:       a1Types.String,
				PrimaryKey: true,
			},
			"_todoId": a1.Field{
				Type:      a1Types.String,
				DataField: "todo_id",
			},
			"_tagName": a1.Field{
				Type:      a1Types.String,
				DataField: "tag_name",
			},

			"tag": a1.LinkedField{
				LinkedModel: "Tag",
				Type:        a1.ManyToOne,
				Field:       "_tagName",
				LinkedField: "_name",
			},
			"todo": a1.LinkedField{
				LinkedModel: "Todo",
				Type:        a1.ManyToOne,
				Field:       "_todoId",
				LinkedField: "_id",
			},
		},
	},
}

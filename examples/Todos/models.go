package main

import (
	"fmt"

	a1 "github.com/aklinker1/a1/pkg"
	a1Types "github.com/aklinker1/a1/pkg/types"
)

func postgreSQLTable(tableName string) a1.DataLoaderConfig {
	return a1.DataLoaderConfig{
		Source: "PostgreSQL",
		Group:  "users",
	}
}

func appendMetadataFields(fieldMap a1.FieldMap) a1.FieldMap {
	fieldMap["createdAt"] = a1.Field{
		Type:      a1Types.Date,
		DataField: "created_at",
	}
	fieldMap["_createdBy"] = a1.Field{
		Type:      a1Types.Date,
		DataField: "created_by",
	}
	fieldMap["updatedAt"] = a1.Field{
		Type:      a1Types.Date,
		DataField: "updated_at",
	}
	fieldMap["_updatedBy"] = a1.Field{
		Type:      a1Types.Date,
		DataField: "updated_by",
	}
	fieldMap["deletedAt"] = a1.Field{
		Type:      a1Types.Date,
		DataField: "deleted_at",
	}
	fieldMap["_deletedBy"] = a1.Field{
		Type:      a1Types.Date,
		DataField: "deleted_by",
	}
	return fieldMap
}

var models = a1.ModelMap{
	"User": a1.Model{
		DataLoader: postgreSQLTable("users"),
		Fields: appendMetadataFields(a1.FieldMap{
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
		}),
	},
	"MyUser": a1.Model{
		Extends: "User",
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
		Fields: appendMetadataFields(a1.FieldMap{
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
			"todoTags": a1.LinkedField{
				LinkedModel: "TodoTag",
				Type:        a1.OneToMany,
				Field:       "_id",
				LinkedField: "_todoId",
			},
		}),
	},
	"Tag": a1.Model{
		DataLoader: postgreSQLTable("tags"),
		Fields: appendMetadataFields(a1.FieldMap{
			"_name": a1.Field{
				Type:       a1Types.String,
				PrimaryKey: true,
			},

			"todoTags": a1.LinkedField{
				LinkedModel: "TodoTag",
				Type:        a1.OneToMany,
				Field:       "_name",
				LinkedField: "_tagName",
			},
		}),
	},
	"TodoTag": a1.Model{
		DataLoader: postgreSQLTable("todo_tags"),
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

var customTypes = a1.CustomTypeMap{
	"Email": a1.CustomType{
		ToJSON: func(value interface{}) interface{} {
			return value.(string)
		},
		FromJSON: func(value interface{}) interface{} {
			return value.(string)
		},
		FromLiteral: func(value a1.ASTValue) interface{} {
			return value.GetValue().(string)
		},
	},
}

var customEnums = a1.EnumMap{
	"Theme": a1.Enum{
		Description: "The UI style of the application",
		Values: a1.EnumValueMap{
			"LIGHT": a1.EnumValue{
				Value: 0,
			},
			"DARK": a1.EnumValue{
				Value: 1,
			},
			"DAY_NIGHT": a1.EnumValue{
				Value:       2,
				Description: "When it is dark out, switch to the dark theme, otherwise be light",
			},
		},
	},
	"Validation": a1.Enum{
		Description: "Whether not not the user has had their email validated",
		Values: a1.EnumValueMap{
			"VERIFIED": a1.EnumValue{
				Value:       0,
				Description: "Once the user has confirmed their email",
			},
			"UNVERIFIED": a1.EnumValue{
				Value:       1,
				Description: "When the user is newly created, and has not verified their email",
			},
			"RECOMMENDED": a1.EnumValue{
				Value:       2,
				Description: "If a user recommends another person to the service",
			},
		},
	},
}

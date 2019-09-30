# Ideal Usage Example

```go
func PostgreSQLTable(tableName string) a1.DataLoaderConfig{
    return a1.DataLoaderConfig{
        Source: "PostgreSQL",
        Group:  "users",
    },
}

func AppendMetadataFields(fieldMap a1.FieldMap) a1.FieldMap {
    fieldMap["createdAt"] = a1.Field{
        Type:      "Date",
        DataField: "created_at",
    },
    fieldMap["_createdBy"] = a1.Field{
        Type:      "Date",
        DataField: "created_by",
    },
    fieldMap["updatedAt"] = a1.Field{
        Type:      "Date",
        DataField: "updated_at",
    },
    fieldMap["_updatedBy"] = a1.Field{
        Type:      "Date",
        DataField: "updated_by",
    },
    fieldMap["deletedAt"] = a1.Field{
        Type:      "Date",
        DataField: "deleted_at",
    },
    fieldMap["_deletedBy"] = a1.Field{
        Type:      "Date",
        DataField: "deleted_by",
    },
    return fieldMap;
}

var models = a1.ModelMap{
    "User": a1.Model{
        dataLoader: PostgreSQLTable("users"),
        Fields: AppendMetadataFields(a1.FieldMap{
            "_id": a1.Field{
                Type:       a1.Types.ID,
                PrimaryKey: true,
            },
            "username": a1.Types.String,
            "email": a1.Field{
                Type:    "Email",
                Graphql: false,
            },
            "passwordHash": a1.Field{
                Type:       a1.Types.String,
                Graphql:    false,
                DataField:  "password_hash",
            },
        }),
    },
    "MyUser": a1.Model{
        Parent: "User",
        Fields: a1.FieldMap{
            "email":      "Email",
            "validation": "Validation",
            "firstName":  a1.Types.String,
            "lastName":   a1.Types.String,

            "fullName": a1.ComputedField{
                Type:     a1.Types.String,
                requires: ["firstName", "lastName"],
                compute: func(data map[string]interface{}) string {
                    return fmt.Sprintf("%s %s", data["firstName"], data["lastName"])
                }
            },

            "todos": a1.LinkedField{
                LinkedModel: "Preferences",
                Type:        a1.LinkTypes.OneToMany,
                Fields:      ["MyUser._id", "Todo._userId"]
            },
            "preferences": a1.LinkedField{
                LinkedModel: "Preferences",
                Type:        a1.LinkTypes.OneToOne,
                Fields:      ["MyUser._id", "Preferences._userId"],
            },
        }
    }
    "Preferences": a1.Model{
        DataLoader: PostgreSQLTable("preferences"),
        Fields: a1.FieldMap{
            "_id": a1.Field{
                Type:       a1.Types.ID,
                PrimaryKey: true,
            },
            "_userId": a1.Field{
                Type: a1.Types.ID,
                DataField: "user_id"
            },
            "theme": "Theme",

            "user": a1.LinkedField{
                LinkedModel: "MyUser",
                Type:        a1.LinkTypes.OneToOne,
                Fields:      ["MyUser._id", "Preferences._userId"]
            },
        },
    },
    "Todo": a1.Model{
        DataLoader: PostgreSQLTable("todos"),
        Fields: AppendMetadataFields(a1.FieldMap{
            "_id": a1.Field{
                Type:       a1.Types.ID,
                PrimaryKey: true,
            },
            "_userId": a1.Field{
                Type:      a1.Types.ID,
                DataField: "user_id"
            },
            "message": a1.Types.String,
            "isCompleted": a1.Field{
                Type:      a1.Types.Boolean,
                DataField: "is_completed",
            },

            "user": a1.LinkedField{
                LinkedModel: "User",
                Type:        a1.LinkTypes.ManyToOne,
                Fields:      ["User._id", "Todo._userId"]
            },
            "todoTags": a1.LinkedField{
                LinkedModel: "TodoTag",
                Type:        a1.LinkTypes.OneToMany,
                Fields:      ["Todo._id", "TodoTag._todoId"],
            }
        }),
    },
    "Tag": a1.Model{
        DataLoader: PostgreSQLTable("tags"),
        Fields: AppendMetadataFields(a1.Fields{
            "_name": a1.Field{
                Type:       a1.Types.String,
                PrimaryKey: true,
            },

            "todoTags": a1.LinkedField{
                LinkedModel: "TodoTag",
                Type:        a1.LinkTypes.OneToMany,
                Fields:      ["Tag._name", "TodoTag._tagName"],
            }
        }),
    },
    "TodoTag": a1.Model{
        DataLoader: PostgreSQLTable("todo_tags"),
        Fields: a1.FieldMap{
            "_id": a1.Field{
                Type:       a1.Types.String,
                PrimaryKey: true,
            },
            "_todoId": a1.Field{
                Type:       a1.Types.String,
                DataField:  "todo_id"
            },
            "_tagName": a1.Field{
                Type:       a1.Types.String,
                DataField:  "tag_name"
            },

            "tag": a1.LinkedField{
                LinkedModel: "Tag",
                Type:        a1.LinkTypes.ManyToOne,
                Fields:      ["Tag._name", "TodoTag._tagName"],
            }
            "todo": a1.LinkedField{
                LinkedModel: "Todo",
                Type:        a1.LinkTypes.ManyToOne,
                Fields:      ["Todo._id", "TodoTag._todoId"],
            }
        },
    },
}

var customTypes = a1.CustomTypeMap{
    "Email": a1.CustomType{
        ToJSON: func(value interface{}) interface{} {
            return ""
        },
        FromJSON: func(value interface{}) interface{} {
            return ""
        },
        FromLiteral: func(value interface{}) interface{} {
            return ""
        },
    }
}

var enums = a1.EnumMap{
    "Theme": a1.IntEnum{
        0: "Light",
        1: "Dark",
        2: "Day/Night",
    },
    "Validation": a1.StringEnum{
        "verified": "Once the user has confirmed their email",
        "new": "When the user is newly created, and has not verified their email",
        "recommended": "If a user recommends another person to the service",
    },
}
```

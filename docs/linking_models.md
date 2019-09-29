# Linking Models

Linking models are how relationships are defined in A1. This page will walk you though the basics of defining a `LinkedField`.

Lets say you are working on a Todo application. Your app has users and todos, and these models might look something like this to start:

```go
var todo = a1.Model{
    Table:      "todos",
    PrimaryKey: "todoId",
    Fields: map[string]a1.Field{
        "todoId": a1.Field{
            Type: "Int",
        },
        "title": a1.Field{
            Type: "String",
        },
    },
}
```

```go
var user = a1.Model{
    Table:      "users",
    PrimaryKey: "userId",
    Fields: map[string]a1.Field{
        "userId": a1.Field{
            Type: "Int",
        },
        "username": a1.Field{
            Type: "String",
        },
    },
}
```

## One To Many

For the example types above, we want to define a relationship between 1 user to many todos. The atual GraphQL queries we want to be able to make look like this:

```graphql
{
    todo(todoId: 1) {
        todoId
        title
        user {
            userId
            username
        }
    }
}
```

```graphql
{
    user(userId: 1) {
        userId
        username
        todos {
            todoId
            title
        }
    }
}
```

To do so, we need to add a field to the `todo` model to keep track of which user the todo belongs to. Lets fall the field `_userId`, but it would be anything.

```go
var todo = a1.Model{
    Table:      "todos",
    PrimaryKey: "todoId",
    Fields: map[string]a1.Field{
        "todoId": a1.Field{
            Type: "Int",
        },
        "title": a1.Field{
            Type: "String",
        },
        // Added field
        "_userId": a1.Field{
            Type: "String",
        },
    },
}
```

If we start up the server like this, the graphql schema will not include the `user`/`todos` fields as part of the todo/user query. To have those fields show up, we need to add a `Linking` argument to our `_userId` field.

```go
"_userId": a1.Field{
    Type: "String",
    Linking: &a1.LinkedField{
        AccessedAs:        "user",
        ModelName:         "User",
        ReverseAccessedAs: "todos",
        LinkedField:       "id",
        Type:              a1.OneToMany,
    },
},
```

Adding the above will add both `Todo.user` and the `User.todos` to the graphql schema. You do not need to add a field on `User` for `todos`, it is done here when `ReverseAccessedAs` is added. If `ReverseAccessedAs` is excluded, it would not create the `User.todos` entry in the graphql schema. Likewise, excluding the Accessed 

Because this is a `OneToMany` relationship, the reverse accessed field `User.todos` will be a list.

## One to One

Creating `OneToOne` relationships is very similar to the `OneToMany`, there being two differences:

1. The reverse accessed field is not a list, but a single object
1. The `LinkedField` argument can go on either objects

So lets say we want our user's to have preferences. I'm building my models for PostgreSQL, so rather than adding another indexed column on my preferences table, I'll add the `LinkedField` to the user model

```go
var preferences = a1.Model{
    Table:      "preferences",
    PrimaryKey: "id",
    Fields: map[string]a1.Field{
        "id": a1.Field{
            Type: "Int",
        },
        "notifications": a1.Field{
            Type: "Boolean",
        },
    },
}
```

```go
var user = a1.Model{
    Table:      "users",
    PrimaryKey: "id",
    Fields: map[string]a1.Field{
        "id": a1.Field{
            Type: "Int",
        },
        "username": a1.Field{
            Type: "String",
        },
        // Added a field for the preference ID and linked it to the "Preferences" model
        "_preferencesId": a1.Field{
            Type: "Int",
            Linking: &a1.LinkedField{
                AccessedAs:        "preferences",
                ModelName:         "Preferences",
                ReverseAccessedAs: "user",
                LinkedField:       "id",
                Type:              a1.OneToOne,
            },
        },
    },
}
```

## Many to Many

Many to many relationships require another model to represent the relationship between the original two models.

Lets say we want to add tags to todos. We have to add another tags model.

```go
var tag = a1.Model{
    Table:      "tags",
    PrimaryKey: "tagId",
    Fields: map[string]a1.Field{
        "tagId": a1.Field{
            Type: "Int",
        },
        "name": a1.Field{
            Type: "String",
        },
    },
}
```

And finally define the many to many relationship

```go
var todoTags = a1.Model{
    Table:      "todo_tags",
    PrimaryKey: "id",
    Fields: map[string]a1.Field{
        "id": a1.Field{
            Type: "Int",
        },
        "_todoId": a1.Field{
            Type: "Int",
            Linking: &a1.LinkedField{
                AccessedAs:        "todo",
                ModelName:         "Todo",
                ReverseAccessedAs: "todos",
                LinkedField:       "id",
                Type:              a1.OneToMany,
            },
        },
        "_tagId": a1.Field{
            Type: "Int",
            Linking: &a1.LinkedField{
                AccessedAs:        "tag",
                ModelName:         "Tag",
                ReverseAccessedAs: "tags",
                LinkedField:       "tagId",
                Type:              a1.OneToMany,
            },
        },
    },
}
```

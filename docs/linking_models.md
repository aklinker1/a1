# <img height="25" src="https://user-images.githubusercontent.com/10101283/66178622-8f14d480-e62b-11e9-8db7-d18cc7885fb3.png"> &ensp; Linking Models

Linking models are how relationships are defined in A1. This page will walk you though the basics of defining a `LinkedField`.

Lets say you are working on a Todo application. Your app has users and todos, and these models might look something like this to start:

```go
var todo = a1.Model{
    // ...
    Fields: map[string]a1.Field{
        "todoId": a1.Field{
            Type:       a1Types.ID,
            PrimaryKey: true,
        },
        "title": a1Types.String,
    },
}
```

```go
var user = a1.Model{
    // ...
    Fields: map[string]a1.Field{
        "userId": a1.Field{
            Type:       a1Types.ID,
            PrimaryKey: true,
        },
        "username": a1Types.String,
    },
}
```

## One To Many

For the example types above, we want to define a relationship between 1 user to many todos. The actual GraphQL queries we want to be able to make look something like this:

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
var user = a1.Model{
    // ...
    Fields: map[string]a1.Field{
        "todoId": a1.Field{
            Type:       a1Types.ID,
            PrimaryKey: true,
        },
        "title": a1Types.String,
        // Added field
        "_userId": a1Types.ID,
    },
}
```

If we start up the server like this, the graphql schema will not include the `user`/`todos` fields as part of the todo/user query. To have those fields show up, we need to add a `LinkedField` to our models

```go
var user = a1.Model{
    // ...
    Fields: map[string]a1.Field{
        "userId": a1.Field{
            Type:       a1Types.ID,
            PrimaryKey: true,
        },
        "username": a1Types.String,

        "todos": a1.LinkedField{
            LinkedModel: "Todo",
            // There is a single user that has many todos
            Type:        a1.OneToMany,
            Field:       "userId",
            LinkedField: "_userId",
        },
    },
}

var todo = a1.Model{
    // ...
    Fields: map[string]a1.Field{
        "todoId": a1.Field{
            Type:       a1Types.ID,
            PrimaryKey: true,
        },
        "title": a1Types.String,

        "_userId": a1Types.ID,
        "user": a1.LinkedField{
            LinkedModel: "User",
            // There are many todos on a single user
            Type:        a1.ManyToOne,
            Field:       "_userId",
            LinkedField: "userId",
        },
    },
}
```

## One to One

Creating `OneToOne` relationships is very similar to the `ManyToOne`, there being two differences:

1. Both sides should use `OneToOne`
2. Both linked fields will return a model, not a list of models

So lets say we want our user's to have preferences. We still have to declare a regular field to act as the connection between the models, but for `OneToOne` relationships, it can go on either model.

That does not mean that it SHOULD go on either model. To make queries faster with whatever data loader you are using, it might make more sense to put it on one model rather than the other. In this example, lets assume the data loader for both the user and preferences is PostgreSQL. Since the majority of the use cases for getting preferences would be to get them from the user, lets add the field there. That way PostgreSQL just has to select the preferences by their primary key for the most use case.

```go
var preferences = a1.Model{
    // ...
    Fields: map[string]a1.Field{
        "preferencesId": a1.Field{
            Type:       a1Types.ID,
            PrimaryKey: true,
        },
        "notifications": a1Types.Boolean,

        "user": a1.LinkedField{
            LinkedModel: "User",
            Type:        a1.OneToOne,
            Field:       "preferencesId",
            LinkedField: "_preferencesId",
        },
    },
}

var user = a1.Model{
    // ...
    Fields: map[string]a1.Field{
        "userId": a1.Field{
            Type:       a1Types.ID,
            PrimaryKey: true,
        },
        "username": a1Types.String,

        "_preferencesId": a1Types.ID,
        "preferences": a1.LinkedField{
            LinkedModel: "Preferences",
            Type:        a1.OneToOne,
            Field:       "_preferencesId",
            LinkedField: "preferencesId",
        },

        "todos": a1.LinkedField{
            LinkedModel: "Todo",
            // There is a single user that has many todos
            Type:        a1.OneToMany,
            Field:       "userId",
            LinkedField: "_userId",
        },
    },
}
```

## Many to Many

Many to many relationships require another model to represent the relationship between the original two models.

Lets say we want to add tags to todos, such that many tags can be on many todos, and vise versa. First, we'll start by adding a tags model.

```go
var tag = a1.Model{
    // ...
    Fields: map[string]a1.Field{
        "tagId": a1.Field{
            Type:       a1Types.ID,
            PrimaryKey: true,
        },
        "name": a1Types.String,
    },
}
```

Next, we'll create the many to many relationship in between model, lets call it `"TodoTag"`. It will need to link to a single todo and a single tag, so lets add the id fields needed here as well.

```go
var todoTag = a1.Model{
    // ...
    Fields: map[string]a1.Field{
        "todoTagId": a1.Field{
            Type:       a1Types.ID,
            PrimaryKey: true,
        },
        "name":    a1Types.String,
        "addedAt":    a1Types.Date,
        "_todoId": a1Types.ID,
        "_tagId":  a1Types.ID,
    },
}
```

And finally, we need to define all the linked fields:

1. `Todo.tags` &rarr; `TodoTag`

    ```go
    var todo = a1.Model{
        // ...
        Fields: map[string]a1.Field{
            "todoId": a1.Field{
                Type:       a1Types.ID,
                PrimaryKey: true,
            },
            "title": a1Types.String,

            "_userId": a1Types.ID,
            "user": a1.LinkedField{
                LinkedModel: "User",
                // There are many todos on a single user
                Type:        a1.ManyToOne,
                Field:       "_userId",
                LinkedField: "userId",
            },

            // Add a list of tags
            "tags": a1.LinkedField{
                LinkedModel: "TodoTag",
                // A single Todo will have many TodoTags
                Type:        a1.OneToMany,
                Field:       "todoId",
                LinkedField: "_todoId",
            }
        },
    }
    ```

2. `Tag.todos` &rarr; `TodoTag`

    ```go
    var tag = a1.Model{
        // ...
        Fields: map[string]a1.Field{
            "tagId": a1.Field{
                Type:       a1Types.ID,
                PrimaryKey: true,
            },
            "name": a1Types.String,

            "todo": a1.LinkedField{
                LinkedModel: "TodoTag",
                // A single tag will have many TodoTags
                Type:        a1.OneToMany,
                Field:       "tagId",
                LinkedField: "_tagId",
            }
        },
    }
    ```

3. `TodoTag.todo` &rarr; `Todo`
4. `TodoTag.tag` &rarr; `Tag`

    ```go
    var todoTag = a1.Model{
        // ...
        Fields: map[string]a1.Field{
            "todoTagId": a1.Field{
                Type:       a1Types.ID,
                PrimaryKey: true,
            },
            "name":    a1Types.String,
            "addedAt":    a1Types.Date,

            "_todoId": a1Types.ID,
            "todo": a1.LinkedField{
                LinkedModel: "Todo",
                // Many TodoTags will belong to one Todo
                Type:        a1.ManyToOne,
                Field:       "_todoId",
                LinkedField: "todoId",
            }

            "_tagId":  a1Types.ID,
            "tag": a1.LinkedField{
                LinkedModel: "Tag",
                // Many TodoTags will belong to one Tag
                Type:        a1.ManyToOne,
                Field:       "_tagId",
                LinkedField: "tagId",
            }
        },
    }
    ```

And there we go! A `ManyToMany` relationship has been setup. You could now do a query to get the tags and when they were added to a todo

```graphql
{
    todo(todoId: 1) {
        title
        tags {
            addedAt
            tag {
                tagId
                name
            }
        }
    }
}
```

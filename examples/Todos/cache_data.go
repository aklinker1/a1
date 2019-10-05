package main

var cachedData = map[string]map[interface{}]map[string]interface{}{
	"users":       users,
	"preferences": preferences,
	"todos":       todos,
	"tags":        tags,
	"todo_tags":   todoTags,
}

var users = map[interface{}]map[string]interface{}{
	1: map[string]interface{}{
		"id":         1,
		"username":   "aklinker1",
		"email":      "example1@gmail.com",
		"validation": 0,
		"firstName":  "Aaron",
		"lastName":   "Klinker",
	},
	2: map[string]interface{}{
		"id":         2,
		"username":   "lonewalker",
		"email":      "example2@gmail.com",
		"validation": 1,
		"firstName":  "Isaiah",
		"lastName":   "Walker",
	},
	3: map[string]interface{}{
		"id":         3,
		"username":   "klinker44",
		"email":      "example3@gmail.com",
		"validation": 2,
		"firstName":  "Luke",
		"lastName":   "Klinker",
	},
}

var preferences = map[interface{}]map[string]interface{}{
	11: map[string]interface{}{
		"id":      11,
		"user_id": 1,
		"theme":   0,
	},
	12: map[string]interface{}{
		"id":      12,
		"user_id": 2,
		"theme":   1,
	},
	13: map[string]interface{}{
		"id":      13,
		"user_id": 3,
		"theme":   2,
	},
}

var todos = map[interface{}]map[string]interface{}{
	1: map[string]interface{}{
		"id":           1,
		"message":      "Todo 1",
		"user_id":      1,
		"is_completed": false,
	},
	2: map[string]interface{}{
		"id":           2,
		"message":      "Todo 2",
		"user_id":      1,
		"is_completed": true,
	},
	3: map[string]interface{}{
		"id":           3,
		"message":      "Todo 3",
		"user_id":      2,
		"is_completed": false,
	},
	4: map[string]interface{}{
		"id":           4,
		"message":      "Todo 4",
		"user_id":      -1,
		"is_completed": true,
	},
}

var todoTags = map[interface{}]map[string]interface{}{
	1: map[string]interface{}{
		"id":       1,
		"todo_id":  1,
		"tag_name": "Tag #1",
		"addedAt":  "Jan 1, 1995",
	},
	2: map[string]interface{}{
		"id":       2,
		"todo_id":  1,
		"tag_name": "Tag #2",
		"addedAt":  "Jan 2, 1995",
	},
	3: map[string]interface{}{
		"id":       3,
		"todo_id":  1,
		"tag_name": "Tag #3",
		"addedAt":  "Jan 3, 1995",
	},
	4: map[string]interface{}{
		"id":       4,
		"todo_id":  2,
		"tag_name": "Tag #3",
		"addedAt":  "Jan 4, 1995",
	},
	5: map[string]interface{}{
		"id":       5,
		"todo_id":  2,
		"tag_name": "Tag #2",
		"addedAt":  "Jan 5, 1995",
	},
	6: map[string]interface{}{
		"id":       6,
		"todo_id":  3,
		"tag_name": "Tag #1",
		"addedAt":  "Jan 6, 1995",
	},
}

var tags = map[interface{}]map[string]interface{}{
	"Tag #1": map[string]interface{}{
		"name": "Tag #1",
	},
	"Tag #2": map[string]interface{}{
		"name": "Tag #2",
	},
	"Tag #3": map[string]interface{}{
		"name": "Tag #3",
	},
}

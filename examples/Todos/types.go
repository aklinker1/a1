package main

import (
	"fmt"

	a1 "github.com/aklinker1/a1/pkg"
)

var customTypes = a1.CustomTypeMap{
	"Email": a1.CustomType{
		ToJSON: func(value interface{}) interface{} {
			if v, ok := value.(*string); ok {
				if v == nil {
					return nil
				}
				return *v
			}
			return fmt.Sprintf("%v", value)
		},
		FromJSON: func(value interface{}) interface{} {
			if v, ok := value.(*string); ok {
				if v == nil {
					fmt.Println("1")
					return nil
				}
				fmt.Println("2")
				return *v
			}
			return fmt.Sprintf("%v", value)
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

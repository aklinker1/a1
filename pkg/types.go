package pkg

import (
	"strconv"

	"github.com/graphql-go/graphql/language/ast"
)

var extendedTypes = CustomTypeMap{
	Int64:   int64Type,
	Float64: float64Type,
}

var int64Type = CustomType{
	Description: "The `Int64` scalar type represents non-fractional numeric values. `Int64` can represent values between -(2^63) and 2^63 - 1.",
	ToJSON:      coerceInt64,
	FromJSON:    coerceInt64,
	FromLiteral: func(valueAST ASTValue) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			if int64Value, err := strconv.ParseInt(valueAST.Value, 10, 64); err == nil {
				return int64Value
			}
		}
		return nil
	},
}

func coerceInt64(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	switch value := value.(type) {
	case int:
		return int64(value)
	case int8:
		return int64(value)
	case int16:
		return int64(value)
	case int32:
		return int64(value)
	case int64:
		return value

	case *int:
		return int64(*value)
	case *int8:
		return int64(*value)
	case *int16:
		return int64(*value)
	case *int32:
		return int64(*value)
	case *int64:
		return int64(*value)

	case float32:
		return int64(value)
	case float64:
		return int64(value)

	case *float32:
		return int64(*value)
	case *float64:
		return int64(*value)

	case string:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil
		}
		return val
	case *string:
		val, err := strconv.ParseInt(*value, 10, 64)
		if err != nil {
			return nil
		}
		return val
	}
	return nil
}

var float64Type = CustomType{
	Description: "The `Float64` scalar type represents non-fractional signed whole numeric " +
		"values. Float64 can represent values between -(2^63) and 2^63 - 1.",
	ToJSON:   coerceFloat64,
	FromJSON: coerceFloat64,
	FromLiteral: func(valueAST ASTValue) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			if int64Value, err := strconv.ParseInt(valueAST.Value, 10, 64); err == nil {
				return int64Value
			}
		}
		return nil
	},
}

func coerceFloat64(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	switch value := value.(type) {
	case int:
		return float64(value)
	case int8:
		return float64(value)
	case int16:
		return float64(value)
	case int32:
		return float64(value)
	case int64:
		return value

	case *int:
		return float64(*value)
	case *int8:
		return float64(*value)
	case *int16:
		return float64(*value)
	case *int32:
		return float64(*value)
	case *int64:
		return float64(*value)

	case float32:
		return float64(value)
	case float64:
		return float64(value)

	case *float32:
		return float64(*value)
	case *float64:
		return float64(*value)

	case string:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil
		}
		return val
	case *string:
		val, err := strconv.ParseFloat(*value, 64)
		if err != nil {
			return nil
		}
		return val
	}
	return nil
}

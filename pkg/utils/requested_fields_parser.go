package utils

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// ParseRequestedFields returns the forward propagated selection from the graphql query
func ParseRequestedFields(params graphql.ResolveParams) (tableFieldsMap map[string]interface{}, err error) {
	fieldASTs := params.Info.FieldASTs
	if len(fieldASTs) == 0 {
		return map[string]interface{}{}, nil
	}
	return selectedFieldsFromSelections(params, fieldASTs[0].SelectionSet.Selections)
}

func selectedFieldsFromSelections(params graphql.ResolveParams, selections []ast.Selection) (map[string]interface{}, error) {
	selected := map[string]interface{}{}
	for _, s := range selections {
		switch t := s.(type) {
		case *ast.Field:
			field := s.(*ast.Field)
			fieldName := field.Name.Value
			if field.SelectionSet == nil {
				selected[fieldName] = true
				break
			}
			sel, err := selectedFieldsFromSelections(params, field.SelectionSet.Selections)
			if err != nil {
				return nil, err
			}
			selected[fieldName] = sel
		case *ast.FragmentSpread:
			fieldName := s.(*ast.FragmentSpread).Name.Value
			frag, ok := params.Info.Fragments[fieldName]
			if !ok {
				return nil, fmt.Errorf("Not fragment found while parsing the data: %v", fieldName)
			}
			sel, err := selectedFieldsFromSelections(params, frag.GetSelectionSet().Selections)
			if err != nil {
				return nil, err
			}
			selected[fieldName] = sel
		default:
			return nil, fmt.Errorf("Unsupported selection type: %v", t)
		}
	}
	return selected, nil
}

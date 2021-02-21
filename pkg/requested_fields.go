package pkg

import (
	"fmt"

	utils "github.com/aklinker1/a1/pkg/utils"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// GetRequestedFields returns the forward propagated selection from the graphql query
func GetRequestedFields(serverConfig *FinalServerConfig, model *FinalModel, params graphql.ResolveParams) (RequestedFieldMap, error) {
	fieldNameMap, err := parseRequestedFieldNameMap(params)
	if err != nil {
		return nil, err
	}

	return createRequestedFieldMap(serverConfig, model, fieldNameMap), nil
}

type parseFieldsHelper struct {
	FieldName string
	Value     interface{}
}

func createRequestedFieldMap(serverConfig *FinalServerConfig, model *FinalModel, fieldNames DataMap) RequestedFieldMap {
	queue := []parseFieldsHelper{}
	for fieldName, value := range fieldNames {
		queue = append(queue, parseFieldsHelper{
			FieldName: fieldName,
			Value:     value,
		})
	}

	fieldMap := RequestedFieldMap{}
	for _, item := range queue {
		fieldName := item.FieldName
		field := model.Fields[fieldName]
		value := item.Value
		requestedField := RequestedField{
			WasRequired: false,
			Field:       field,
		}

		switch field.(type) {
		case *FinalVirtualField:
			virtualField := field.(*FinalVirtualField)
			for _, requiredFieldName := range virtualField.RequiredFields {
				if _, exists := fieldNames[requiredFieldName]; !exists {
					// TODO - Make recursive to handle linked object requirements, or just add the fields to the fieldNames object?
					fieldMap[requiredFieldName] = RequestedField{
						WasRequired: true,
						Field:       model.Fields[requiredFieldName],
						InnerFields: nil,
						Model:       nil,
					}
				}
			}
		case *FinalLinkedField:
			nextFieldNames := value.(DataMap)
			linkedField := field.(*FinalLinkedField)
			requestedField.Model = linkedField.LinkedModel
			requestedField.InnerFields = createRequestedFieldMap(serverConfig, linkedField.LinkedModel, nextFieldNames)
			requiredFieldName := linkedField.Field
			if _, exists := fieldNames[requiredFieldName]; !exists {
				fieldMap[requiredFieldName] = RequestedField{
					WasRequired: true,
					Field:       model.Fields[requiredFieldName],
					InnerFields: nil,
					Model:       nil,
				}
			}
		}

		fieldMap[fieldName] = requestedField
	}
	return fieldMap
}

func parseRequestedFieldNameMap(params graphql.ResolveParams) (DataMap, error) {
	fieldASTs := params.Info.FieldASTs
	if len(fieldASTs) == 0 {
		return DataMap{}, nil
	}
	return selectedFieldsFromSelections(params, fieldASTs[0].SelectionSet.Selections)
}

func selectedFieldsFromSelections(params graphql.ResolveParams, selections []ast.Selection) (DataMap, error) {
	selected := DataMap{}
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

func printRequestedFields(requestedFields RequestedFieldMap, indentSpaces int) {
	indent := ""
	for i := 0; i < indentSpaces; i++ {
		indent += " "
	}
	for fieldName, field := range requestedFields {
		requiredString := ""
		if field.WasRequired {
			requiredString = " *"
		}
		utils.Log("%s- %s%s", indent, fieldName, requiredString)
		if field.InnerFields != nil {
			printRequestedFields(field.InnerFields.(RequestedFieldMap), indentSpaces+2)
		}
	}
}

package pkg

import utils "github.com/aklinker1/a1/pkg/utils"

// GetOneQuery -
func GetOneQuery(serverConfig *FinalServerConfig, model *FinalModel) Resolvable {
	return Resolvable{
		Name:  utils.LowerFirstChar(model.Name),
		Model: model,
		Arguments: []Argument{
			Argument{
				Name: model.PrimaryKey,
				Type: "Int",
			},
		},
		Resolver: func(args DataMap, requestedFields RequestedFieldMap) (interface{}, error) {
			return getOneModel(serverConfig, model, args, requestedFields)
		},
	}
}

// GetMultipleQuery -
func GetMultipleQuery(serverConfig *FinalServerConfig, model *FinalModel) Resolvable {
	return Resolvable{
		Name:         "list" + utils.AddS(model.Name),
		Model:        model,
		ResturnsList: true,
		Resolver: func(args DataMap, requestedFields RequestedFieldMap) (interface{}, error) {
			return getMultipleModels(serverConfig, model, args, requestedFields)
		},
	}
}

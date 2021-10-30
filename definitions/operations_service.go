package definitions

var OperationsService = []Operation{
	OperationServiceCreate,
	OperationServiceDelete,
	OperationServiceGet,
	OperationServiceList,
}

var OperationServiceCreate = Operation{
	Name:      "create",
	Namespace: NamespaceService,
	Params: []Field{
		getField("name"),
	},
	Results: []Field{
		getField("store"),
	},
	Description: "will create a new storager instance.",
}

var OperationServiceDelete = Operation{
	Name:      "delete",
	Namespace: NamespaceService,
	Params: []Field{
		getField("name"),
	},
	Description: "will delete a storager instance.",
}

var OperationServiceGet = Operation{
	Name:      "get",
	Namespace: NamespaceService,
	Params: []Field{
		getField("name"),
	},
	Results: []Field{
		getField("store"),
	},
	Description: "will get a valid storager instance for service.",
}

var OperationServiceList = Operation{
	Name:      "list",
	Namespace: NamespaceService,
	Results: []Field{
		getField("sti"),
	},
	Description: "will list all storager instances under this service.",
}

func init() {
	for k := range OperationsService {
		OperationsService[k].Params = append(OperationsService[k].Params, getField("pairs"))
		if !OperationsService[k].Local {
			OperationsService[k].Results = append(OperationsService[k].Results, getField("err"))
		}
	}
}

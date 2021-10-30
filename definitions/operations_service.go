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
		FieldMap["name"],
	},
	Results: []Field{
		FieldMap["store"],
	},
	Description: "will create a new storager instance.",
}

var OperationServiceDelete = Operation{
	Name:      "delete",
	Namespace: NamespaceService,
	Params: []Field{
		FieldMap["name"],
	},
	Description: "will delete a storager instance.",
}

var OperationServiceGet = Operation{
	Name:      "get",
	Namespace: NamespaceService,
	Params: []Field{
		FieldMap["name"],
	},
	Results: []Field{
		FieldMap["store"],
	},
	Description: "will get a valid storager instance for service.",
}

var OperationServiceList = Operation{
	Name:      "list",
	Namespace: NamespaceService,
	Results: []Field{
		FieldMap["sti"],
	},
	Description: "will list all storager instances under this service.",
}

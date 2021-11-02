package definitions

import (
	"fmt"
	"github.com/Xuanwo/gg"
	"github.com/Xuanwo/templateutils"
	log "github.com/sirupsen/logrus"
)

type genNamespace struct {
	g *gg.Generator
}

func GenerateNamespace(path string) {
	gf := &genNamespace{
		g: gg.New(),
	}

	f := gf.g.NewGroup()
	f.AddLineComment("Code generated by go generate cmd/definitions; DO NOT EDIT.")
	f.AddPackage("definitions")
	f.NewImport().
		AddPath("go.beyondstorage.io/v5/types")

	gf.generateNamespace(NamespaceService, OperationsService)
	gf.generateNamespace(NamespaceStorage, OperationsStorage)

	err := gf.g.WriteFile(path)
	if err != nil {
		log.Fatalf("generate to %s: %v", path, err)
	}
}

func (gf *genNamespace) generateNamespace(nsName string, ops []Operation) {
	f := gf.g.NewGroup()

	nsNameP := templateutils.ToPascal(nsName)

	structName := fmt.Sprintf("%s", nsNameP)
	sf := f.NewStruct(structName)
	sf.AddLine()
	sf.AddField("Features", "types."+nsNameP+"Features")
	sf.AddLine()
	for _, op := range ops {
		sf.AddField(templateutils.ToPascal(op.Name), "[]Pair")
	}
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func format(data *Data) {
	// Generate pairs
	hf := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(data.pairSpec, hf.Body())

	formatBody(hf.Body())

	content := hclwrite.Format(hf.Bytes())
	err := ioutil.WriteFile(pairPath, content, 0644)
	if err != nil {
		log.Fatalf("format: %v", err)
	}

	// Generate metadata
	hf = hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(data.metaSpec, hf.Body())

	formatBody(hf.Body())

	content = hclwrite.Format(hf.Bytes())
	err = ioutil.WriteFile(metadataPath, content, 0644)
	if err != nil {
		log.Fatalf("format: %v", err)
	}

	// Generate services
	for _, v := range data.serviceSpec {
		filePath := fmt.Sprintf("services/%s.hcl", v.Name)

		hf = hclwrite.NewEmptyFile()
		gohcl.EncodeIntoBody(v, hf.Body())

		formatBody(hf.Body())

		content = hclwrite.Format(hf.Bytes())
		err = ioutil.WriteFile(filePath, content, 0644)
		if err != nil {
			log.Fatalf("format: %v", err)
		}
	}
}

func formatBody(body *hclwrite.Body) {
	for k, v := range body.Attributes() {
		if isAttrEmpty(v) {
			body.RemoveAttribute(k)
			continue
		}
	}
	for _, v := range body.Blocks() {
		if isBlockEmpty(v) {
			body.RemoveBlock(v)
			continue
		}

		formatBody(v.Body())
	}
}

func isBlockEmpty(block *hclwrite.Block) bool {
	if len(block.Body().Blocks()) > 0 {
		return false
	}

	attrs := block.Body().Attributes()
	for _, v := range attrs {
		if !isAttrEmpty(v) {
			return false
		}
	}
	return true
}

func isAttrEmpty(attr *hclwrite.Attribute) bool {
	tokens := attr.Expr().BuildTokens(make([]*hclwrite.Token, 0))

	// xxx = null
	if len(tokens) == 1 && tokens[0].Type == hclsyntax.TokenIdent {
		return true
	}
	// xxx = ""
	if len(tokens) == 2 &&
		tokens[0].Type == hclsyntax.TokenOQuote &&
		tokens[1].Type == hclsyntax.TokenCQuote {
		return true
	}
	// xxx = []
	if len(tokens) == 2 &&
		tokens[0].Type == hclsyntax.TokenOBrack &&
		tokens[1].Type == hclsyntax.TokenCBrack {
		return true
	}
	return false
}

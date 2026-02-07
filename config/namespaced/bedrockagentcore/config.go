// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package bedrockagentcore

import (
	"context"
	"fmt"
	"slices"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/crossplane/upjet/v2/pkg/types/name"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func Configure(p *config.Provider) { //nolint:gocyclo
	//p.AddResourceConfigurator("aws_bedrockagentcore_agent_runtime", func(r *config.Resource) {
	//	r.AddSingletonListConversion("agent_runtime_artifact", "agentRuntimeArtifact")
	//	r.AddSingletonListConversion("agent_runtime_artifact[*].code_configuration", "agentRuntimeArtifact[*].codeConfiguration")
	//	r.AddSingletonListConversion("agent_runtime_artifact[*].code_configuration[*].code", "agentRuntimeArtifact[*].codeConfiguration[*].code")
	//	r.AddSingletonListConversion("agent_runtime_artifact[*].code_configuration[*].code[*].s3", "agentRuntimeArtifact[*].codeConfiguration[*].code[*].s3")
	//	r.AddSingletonListConversion("agent_runtime_artifact[*].container_configuration", "agentRuntimeArtifact[*].containerConfiguration")
	//	r.AddSingletonListConversion("authorizer_configuration", "authorizerConfiguration")
	//	r.AddSingletonListConversion("authorizer_configuration[*].custom_jwt_authorizer", "authorizerConfiguration[*].customJwtAuthorizer")
	//	r.AddSingletonListConversion("network_configuration", "networkConfiguration")
	//	r.AddSingletonListConversion("network_configuration[*].network_mode_config", "networkConfiguration[*].networkModeConfig")
	//	r.AddSingletonListConversion("protocol_configuration", "protocolConfiguration")
	//	r.AddSingletonListConversion("request_header_configuration", "requestHeaderConfiguration")
	//
	//})
	p.AddResourceConfigurator("aws_bedrockagentcore_agent_runtime", func(r *config.Resource) {
		schResp := &resource.SchemaResponse{}
		r.TerraformPluginFrameworkResource.Schema(context.TODO(), resource.SchemaRequest{}, schResp)
		sch := schResp.Schema
		traverseSchema(sch)
		for tp, crdP := range slc {
			fmt.Println("/////// Adding Framework Singleton-list conversion TP:", tp, "CRDP:", crdP)
			r.AddSingletonListConversion(tp, crdP)
		}
	})

}

var slc = make(map[string]string)

func traverseSchema(s schema.Schema) {
	// The Schema has an Attributes map

	for attrName, attr := range s.Attributes {
		traverseAttribute(attrName, name.NewFromSnake(attrName).LowerCamelComputed, attr)
	}

	// The Schema also has Blocks
	for blockName, block := range s.Blocks {
		traverseBlock(blockName, name.NewFromSnake(blockName).LowerCamelComputed, block)
	}
}

func traverseAttribute(tfPath string, crdPath string, attr schema.Attribute) {
	switch a := attr.(type) {
	case schema.StringAttribute:
		// Handle string attribute
		fmt.Printf("String attribute: %s\n", tfPath)

	case schema.Int64Attribute:
		fmt.Printf("Int64 attribute: %s\n", tfPath)

	case schema.BoolAttribute:
		fmt.Printf("Bool attribute: %s\n", tfPath)

	case schema.ListAttribute:
		fmt.Printf("List attribute: %s\n", tfPath)
		if slices.ContainsFunc(a.Validators, func(validator validator.List) bool {
			return validator.Description(context.TODO()) == "list must contain at most 1 elements"
		}) {
			slc[tfPath] = crdPath
		}
	// Access element type: a.ElementType
	case schema.SetAttribute:
		fmt.Printf("Set attribute: %s\n", tfPath)
		if slices.ContainsFunc(a.Validators, func(validator validator.Set) bool {
			return validator.Description(context.TODO()) == "set must contain at most 1 elements"
		}) {
			slc[tfPath] = crdPath
		}

	case schema.MapAttribute:
		fmt.Printf("Map attribute: %s\n", tfPath)
		// Access element type: a.ElementType

	case schema.ObjectAttribute:
		fmt.Printf("Object attribute: %s\n", tfPath)
		// Access attribute types: a.AttributeTypes

	case schema.SingleNestedAttribute:
		fmt.Printf("SingleNested attribute: %s\n", tfPath)
		// Recursively traverse nested attributes
		for nestedName, nestedAttr := range a.Attributes {
			traverseAttribute(tfPath+"[*]."+nestedName, crdPath+"[*]."+name.NewFromSnake(nestedName).LowerCamelComputed, nestedAttr)
		}

	case schema.ListNestedAttribute:
		fmt.Printf("ListNested attribute: %s\n", tfPath)
		if slices.ContainsFunc(a.Validators, func(validator validator.List) bool {
			return validator.Description(context.TODO()) == "list must contain at most 1 elements"
		}) {
			slc[tfPath] = crdPath
		}
		// Recursively traverse nested attributes
		for nestedName, nestedAttr := range a.NestedObject.Attributes {
			traverseAttribute(tfPath+"[*]."+nestedName, crdPath+"[*]."+name.NewFromSnake(nestedName).LowerCamelComputed, nestedAttr)
		}

	case schema.MapNestedAttribute:
		fmt.Printf("MapNested attribute: %s\n", tfPath)
		// Similar to ListNested
		for nestedName, nestedAttr := range a.NestedObject.Attributes {
			traverseAttribute(tfPath+"[*]."+nestedName, crdPath+"[*]."+name.NewFromSnake(nestedName).LowerCamelComputed, nestedAttr)
		}

	case schema.SetNestedAttribute:
		fmt.Printf("SetNested attribute: %s\n", tfPath)
		if slices.ContainsFunc(a.Validators, func(validator validator.Set) bool {
			return validator.Description(context.TODO()) == "set must contain at most 1 elements"
		}) {
			slc[tfPath] = crdPath
		}
		// Similar to ListNested
		for nestedName, nestedAttr := range a.NestedObject.Attributes {
			traverseAttribute(tfPath+"[*]."+nestedName, crdPath+"[*]."+name.NewFromSnake(nestedName).LowerCamelComputed, nestedAttr)
		}

	}
}

func traverseBlock(tfPath, crdPath string, block schema.Block) {
	switch b := block.(type) {
	case schema.ListNestedBlock:
		if slices.ContainsFunc(b.Validators, func(validator validator.List) bool {
			return validator.Description(context.TODO()) == "list must contain at most 1 elements"
		}) {
			slc[tfPath] = crdPath
		}
		for attrName, attr := range b.NestedObject.Attributes {
			traverseAttribute(tfPath+"[*]."+attrName, crdPath+"[*]."+name.NewFromSnake(attrName).LowerCamelComputed, attr)
		}
		for blockName, bl := range b.NestedObject.Blocks {
			traverseBlock(tfPath+"[*]."+blockName, crdPath+"[*]."+name.NewFromSnake(blockName).LowerCamelComputed, bl)
		}

	case schema.SetNestedBlock:
		if slices.ContainsFunc(b.Validators, func(validator validator.Set) bool {
			return validator.Description(context.TODO()) == "set must contain at most 1 elements"
		}) {
			slc[tfPath] = crdPath
		}
		for attrName, attr := range b.NestedObject.Attributes {
			traverseAttribute(tfPath+"[*]."+attrName, crdPath+"[*]."+name.NewFromSnake(attrName).LowerCamelComputed, attr)
		}
		for blockName, bl := range b.NestedObject.Blocks {
			traverseBlock(tfPath+"[*]."+blockName, crdPath+"[*]."+name.NewFromSnake(blockName).LowerCamelComputed, bl)
		}

	case schema.SingleNestedBlock:
		for attrName, attr := range b.Attributes {
			traverseAttribute(tfPath+"[*]."+attrName, crdPath+"[*]."+name.NewFromSnake(attrName).LowerCamelComputed, attr)
		}
		for blockName, bl := range b.Blocks {
			traverseBlock(tfPath+"[*]."+blockName, crdPath+"[*]."+name.NewFromSnake(blockName).LowerCamelComputed, bl)
		}
	}
}

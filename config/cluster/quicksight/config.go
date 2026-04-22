// SPDX-FileCopyrightText: 2025 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package quicksight

import (
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/upbound/provider-aws/v2/config/cluster/common"
)

// Configure adds configurations for the quicksight group.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("aws_quicksight_dashboard", func(r *config.Resource) {
		delete(r.TerraformResource.Schema, "definition")
		l := r.TFListConversionPaths()
		for _, e := range l {
			if strings.HasPrefix(e, "definition[*].") {
				r.RemoveSingletonListConversion(e)
			}
		}
		// r.MetaResource.ArgumentDocs["definition_json"] = "A raw JSON string used to define the dashboard structure. When this field is used, Crossplane cannot observe changes in the configuration through the AWS API; therefore, drift detection cannot be performed. Refer to the AWS documentation for the expected JSON structure: https://docs.aws.amazon.com/quicksight/latest/APIReference/API_CreateDashboard.html"
		r.MetaResource.Description = "Creates a QuickSight Dashboard resource. The 'definition' field is not supported due to Kubernetes CRD size limitations with deeply nested fields." // Please use the 'definitionJson' field to define the dashboard structure."
		r.References["source_entity.source_template.data_set_references.data_set_arn"] = config.Reference{
			TerraformName: "aws_quicksight_data_set",
			Extractor:     common.PathARNExtractor,
		}
		r.References["source_entity.source_template.arn"] = config.Reference{
			TerraformName: "aws_quicksight_template",
			Extractor:     common.PathARNExtractor,
		}
		r.References["theme_arn"] = config.Reference{
			TerraformName: "aws_quicksight_theme",
			Extractor:     common.PathARNExtractor,
		}
	})

	p.AddResourceConfigurator("aws_quicksight_analysis", func(r *config.Resource) {
		delete(r.TerraformResource.Schema, "definition")
		l := r.TFListConversionPaths()
		for _, e := range l {
			if strings.HasPrefix(e, "definition[*].") {
				r.RemoveSingletonListConversion(e)
			}
		}
		r.MetaResource.Description = "Creates a QuickSight Analysis resource. The 'definition' field is not supported due to Kubernetes CRD size limitations with deeply nested fields."
		r.References["source_entity.source_template.data_set_references.data_set_arn"] = config.Reference{
			TerraformName: "aws_quicksight_data_set",
			Extractor:     common.PathARNExtractor,
		}
		r.References["source_entity.source_template.arn"] = config.Reference{
			TerraformName: "aws_quicksight_template",
			Extractor:     common.PathARNExtractor,
		}
		r.References["theme_arn"] = config.Reference{
			TerraformName: "aws_quicksight_theme",
			Extractor:     common.PathARNExtractor,
		}
	})

	p.AddResourceConfigurator("aws_quicksight_vpc_connection", func(r *config.Resource) {
		r.References["role_arn"] = config.Reference{
			TerraformName: "aws_iam_role",
			Extractor:     common.PathARNExtractor,
		}
		r.References["security_group_ids"] = config.Reference{
			TerraformName: "aws_security_group",
		}
		r.References["subnet_ids"] = config.Reference{
			TerraformName: "aws_subnet",
		}
	})

	p.AddResourceConfigurator("aws_quicksight_data_source", func(r *config.Resource) {
		r.References["vpc_connection_properties.vpc_connection_arn"] = config.Reference{
			TerraformName: "aws_quicksight_vpc_connection",
			Extractor:     common.PathARNExtractor,
		}
		r.References["credentials.copy_source_arn"] = config.Reference{
			TerraformName: "aws_quicksight_data_source",
			Extractor:     common.PathARNExtractor,
		}
		r.References["credentials.secret_arn"] = config.Reference{
			TerraformName: "aws_secretsmanager_secret",
			Extractor:     common.PathARNExtractor,
		}
		r.References["parameters.s3.role_arn"] = config.Reference{
			TerraformName: "aws_iam_role",
			Extractor:     common.PathARNExtractor,
		}
	})

	p.AddResourceConfigurator("aws_quicksight_data_set", func(r *config.Resource) {
		r.References["physical_table_map.custom_sql.data_source_arn"] = config.Reference{
			TerraformName: "aws_quicksight_data_source",
			Extractor:     common.PathARNExtractor,
		}
		r.References["physical_table_map.relational_table.data_source_arn"] = config.Reference{
			TerraformName: "aws_quicksight_data_source",
			Extractor:     common.PathARNExtractor,
		}
		r.References["physical_table_map.s3_source.data_source_arn"] = config.Reference{
			TerraformName: "aws_quicksight_data_source",
			Extractor:     common.PathARNExtractor,
		}
		r.References["logical_table_map.source.data_set_arn"] = config.Reference{
			TerraformName: "aws_quicksight_data_set",
			Extractor:     common.PathARNExtractor,
		}
		r.References["row_level_permission_data_set.arn"] = config.Reference{
			TerraformName: "aws_quicksight_data_set",
			Extractor:     common.PathARNExtractor,
		}
	})

	p.AddResourceConfigurator("aws_quicksight_account_subscription", func(r *config.Resource) {
		r.References["directory_id"] = config.Reference{
			TerraformName: "aws_directory_service_directory",
		}
	})

	p.AddResourceConfigurator("aws_quicksight_group_membership", func(r *config.Resource) {
		r.References["group_name"] = config.Reference{
			TerraformName: "aws_quicksight_group",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("group_name",false)`,
		}
	})

	p.AddResourceConfigurator("aws_quicksight_template", func(r *config.Resource) {
		delete(r.TerraformResource.Schema, "definition")
		l := r.TFListConversionPaths()
		for _, e := range l {
			if strings.HasPrefix(e, "definition[*].") {
				r.RemoveSingletonListConversion(e)
			}
		}
		r.MetaResource.Description = "Creates a QuickSight Template resource. The 'definition' field is not supported due to Kubernetes CRD size limitations with deeply nested fields."
		r.References["source_entity.source_analysis.arn"] = config.Reference{
			TerraformName: "aws_quicksight_analysis",
			Extractor:     common.PathARNExtractor,
		}
		r.References["source_entity.source_analysis.data_set_references.data_set_arn"] = config.Reference{
			TerraformName: "aws_quicksight_data_set",
			Extractor:     common.PathARNExtractor,
		}
		r.References["source_entity.source_template.arn"] = config.Reference{
			TerraformName: "aws_quicksight_template",
			Extractor:     common.PathARNExtractor,
		}
	})

	p.AddResourceConfigurator("aws_quicksight_template_alias", func(r *config.Resource) {
		r.References["template_id"] = config.Reference{
			TerraformName: "aws_quicksight_template",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("template_id",false)`,
		}
	})

	p.AddResourceConfigurator("aws_quicksight_theme", func(r *config.Resource) {
		r.References["base_theme_id"] = config.Reference{
			TerraformName: "aws_quicksight_theme",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("theme_id",false)`,
		}
	})

	p.AddResourceConfigurator("aws_quicksight_folder", func(r *config.Resource) {
		r.References["parent_folder_arn"] = config.Reference{
			TerraformName: "aws_quicksight_folder",
			Extractor:     common.PathARNExtractor,
		}
	})

	p.AddResourceConfigurator("aws_quicksight_folder_membership", func(r *config.Resource) {
		r.References["folder_id"] = config.Reference{
			TerraformName: "aws_quicksight_folder",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("folder_id",false)`,
		}
	})

	p.AddResourceConfigurator("aws_quicksight_iam_policy_assignment", func(r *config.Resource) {
		r.References["policy_arn"] = config.Reference{
			TerraformName: "aws_iam_policy",
			Extractor:     common.PathARNExtractor,
		}
	})

	p.AddResourceConfigurator("aws_quicksight_ingestion", func(r *config.Resource) {
		r.References["data_set_id"] = config.Reference{
			TerraformName: "aws_quicksight_data_set",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("data_set_id",false)`,
		}
	})

	p.AddResourceConfigurator("aws_quicksight_key_registration", func(r *config.Resource) {
		r.References["key_registration.key_arn"] = config.Reference{
			TerraformName: "aws_kms_key",
			Extractor:     common.PathARNExtractor,
		}
	})

	p.AddResourceConfigurator("aws_quicksight_refresh_schedule", func(r *config.Resource) {
		r.References["data_set_id"] = config.Reference{
			TerraformName: "aws_quicksight_data_set",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("data_set_id",false)`,
		}
	})

	p.AddResourceConfigurator("aws_quicksight_role_custom_permission", func(r *config.Resource) {
		r.References["custom_permissions_name"] = config.Reference{
			TerraformName: "aws_quicksight_custom_permissions",
		}
	})

	p.AddResourceConfigurator("aws_quicksight_role_membership", func(r *config.Resource) {
		r.References["member_name"] = config.Reference{
			TerraformName: "aws_quicksight_group",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("group_name",false)`,
		}
	})

	p.AddResourceConfigurator("aws_quicksight_user_custom_permission", func(r *config.Resource) {
		r.References["custom_permissions_name"] = config.Reference{
			TerraformName: "aws_quicksight_custom_permissions",
		}
		r.References["user_name"] = config.Reference{
			TerraformName: "aws_quicksight_user",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("user_name",false)`,
		}
	})
}

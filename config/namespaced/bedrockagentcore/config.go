// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package bedrockagentcore

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

func Configure(p *config.Provider) { //nolint:gocyclo
	p.AddResourceConfigurator("aws_bedrockagentcore_agent_runtime", func(r *config.Resource) {
		r.AddSingletonListConversion("agent_runtime_artifact", "agentRuntimeArtifact")
		r.AddSingletonListConversion("agent_runtime_artifact[*].code_configuration", "agentRuntimeArtifact[*].codeConfiguration")
		r.AddSingletonListConversion("agent_runtime_artifact[*].code_configuration[*].code", "agentRuntimeArtifact[*].codeConfiguration[*].code")
		r.AddSingletonListConversion("agent_runtime_artifact[*].code_configuration[*].code[*].s3", "agentRuntimeArtifact[*].codeConfiguration[*].code[*].s3")
		r.AddSingletonListConversion("agent_runtime_artifact[*].container_configuration", "agentRuntimeArtifact[*].containerConfiguration")
		r.AddSingletonListConversion("authorizer_configuration", "authorizerConfiguration")
		r.AddSingletonListConversion("authorizer_configuration[*].custom_jwt_authorizer", "authorizerConfiguration[*].customJwtAuthorizer")
		r.AddSingletonListConversion("network_configuration", "networkConfiguration")
		r.AddSingletonListConversion("network_configuration[*].network_mode_config", "networkConfiguration[*].networkModeConfig")
		r.AddSingletonListConversion("protocol_configuration", "protocolConfiguration")
		r.AddSingletonListConversion("request_header_configuration", "requestHeaderConfiguration")

	})
}

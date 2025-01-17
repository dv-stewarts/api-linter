// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aep0004

import (
	"fmt"

	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/locations"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var resourceTypeName = &lint.MessageRule{
	Name:     lint.NewRuleName(4, "resource-type-name"),
	RuleType: lint.NewRuleType(lint.MustRule),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		return utils.GetResource(m) != nil
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		resource := utils.GetResource(m)
		_, typeName, ok := utils.SplitResourceTypeName(resource.GetType())
		kebabCase := utils.ToKebabCase(typeName)
		if !ok {
			return []lint.Problem{{
				Message:    "Resource type names must be of the form {Service Name}/{Type}.",
				Descriptor: m,
				Location:   locations.MessageResource(m),
			}}
		}
		if kebabCase != typeName {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Type must be kebob-case with alphanumeric characters: %q", kebabCase),
				Descriptor: m,
				Location:   locations.MessageResource(m),
			}}
		}
		return nil
	},
}

// Copyright 2023 Google LLC
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
	"regexp"
	"strings"

	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/locations"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var identifierRegexp = regexp.MustCompile("^{[a-z][-a-z0-9]*[a-z0-9]}$")

var resourceNameComponentsAlternate = &lint.MessageRule{
	Name:     lint.NewRuleName(4, "resource-name-components-alternate"),
	RuleType: lint.NewRuleType(lint.MustRule),
	OnlyIf:   utils.IsResource,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		var problems []lint.Problem
		resource := utils.GetResource(m)
		for _, p := range resource.GetPattern() {
			components := strings.Split(p, "/")
			for i, c := range components {
				identifierExpected := i%2 == 1
				if identifierExpected != isIdentifier(c) {
					problems = append(problems, lint.Problem{
						Message:    fmt.Sprintf("Resource pattern %q must alternate between collection and identifier. %q is not an identifier", p, c),
						Descriptor: m,
						Location:   locations.MessageResource(m),
					})
					break
				}
			}
		}
		return problems
	},
}

func isIdentifier(s string) bool {
	return identifierRegexp.MatchString(s)
}

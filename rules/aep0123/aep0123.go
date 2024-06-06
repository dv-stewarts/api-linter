// Copyright 2019 Google LLC
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

// Package aep0123 contains rules defined in https://aep.dev/123.
package aep0123

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

// AddRules accepts a register function and registers each of
// this AEP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		123,
		duplicateResource,
		resourceAnnotation,
		resourceNameComponentsAlternate,
		resourceNameField,
		resourcePattern,
		resourcePlural,
		resourceReferenceType,
		resourceSingular,
		resourceTypeName,
		resourceVariables,
		resourceDefinitionVariables,
		resourceDefinitionPatterns,
		resourceDefinitionTypeName,
		nameNeverOptional,
	)
}

func isResourceMessage(m *desc.MessageDescriptor) bool {
	// If the parent of this message is a message, it is nested and shoudn't
	// be considered a resource, even if it has a name field.
	_, nested := m.GetParent().(*desc.MessageDescriptor)
	return m.FindFieldByName("name") != nil && !strings.HasSuffix(m.GetName(), "Request") &&
		!strings.HasSuffix(m.GetName(), "Response") && !nested
}

func hasResourceAnnotation(m *desc.MessageDescriptor) bool {
	return utils.GetResource(m) != nil
}

func hasResourceDefinitionAnnotation(f *desc.FileDescriptor) bool {
	return len(utils.GetResourceDefinitions(f)) > 0
}

// getVariables returns a slice of variables declared in the pattern.
//
// For example, a pattern of "publishers/{publisher}/books/{book}" would
// return []string{"publisher", "book"}.
func getVariables(pattern string) []string {
	answer := []string{}
	for _, match := range varRegexp.FindAllStringSubmatch(pattern, -1) {
		answer = append(answer, match[1])
	}
	return answer
}

// getPlainPattern returns the pattern with all variables replaced with "*".
//
// For example, a pattern of "publishers/{publisher}/books/{book}" would
// return "publishers/*/books/*".
func getPlainPattern(pattern string) string {
	return varRegexp.ReplaceAllLiteralString(pattern, "*")
}

// getDesiredPattern returns the expected desired pattern, with errors we
// lint for corrected.
func getDesiredPattern(pattern string) string {
	want := []string{}
	for _, token := range strings.Split(pattern, "/") {
		if strings.HasPrefix(token, "{") && strings.HasSuffix(token, "}") {
			varname := token[1 : len(token)-1]
			want = append(want, fmt.Sprintf("{%s}", strings.TrimSuffix(strcase.SnakeCase(varname), "_id")))
		} else {
			want = append(want, strcase.LowerCamelCase(token))
		}
	}
	return strings.Join(want, "/")
}

var varRegexp = regexp.MustCompile(`\{([^}=]+)}`)
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

package aep0131

import (
	"fmt"

	"bitbucket.org/creachadair/stringset"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var allowedFields = stringset.New(
	"path",       // AEP-131
	"request_id", // AEP-155
	"read_mask",  // AEP-157
	"view",       // AEP-157
)

// Get methods should not have unrecognized fields.
var unknownFields = &lint.FieldRule{
	Name: lint.NewRuleName(131, "request-unknown-fields"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return utils.IsGetRequestMessage(f.GetOwner())
	},
	RuleType: lint.NewRuleType(lint.MustRule),
	LintField: func(field *desc.FieldDescriptor) []lint.Problem {
		if !allowedFields.Contains(field.GetName()) {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Unexpected field: Get RPCs must only contain fields explicitly described in https://aep.dev/131, not %q.",
					string(field.GetName()),
				),
				Descriptor: field,
			}}
		}

		return nil
	},
}
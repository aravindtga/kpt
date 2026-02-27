// Copyright 2026 The kpt Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRenderCondition(t *testing.T) {
	cond := NewRenderCondition(ConditionTrue, ReasonRenderSucceeded, "")
	assert.Equal(t, ConditionTypeRendered, cond.Type)
	assert.Equal(t, ConditionTrue, cond.Status)
	assert.Equal(t, ReasonRenderSucceeded, cond.Reason)
	assert.Empty(t, cond.Message)

	cond = NewRenderCondition(ConditionFalse, ReasonRenderFailed, "set-annotations failed")
	assert.Equal(t, ConditionTypeRendered, cond.Type)
	assert.Equal(t, ConditionFalse, cond.Status)
	assert.Equal(t, ReasonRenderFailed, cond.Reason)
	assert.Equal(t, "set-annotations failed", cond.Message)
}

func TestStatus_SetCondition_Adds(t *testing.T) {
	s := &Status{}
	cond := NewRenderCondition(ConditionTrue, ReasonRenderSucceeded, "")
	s.SetCondition(cond)

	assert.Len(t, s.Conditions, 1)
	assert.Equal(t, ConditionTypeRendered, s.Conditions[0].Type)
	assert.Equal(t, ConditionTrue, s.Conditions[0].Status)
}

func TestStatus_SetCondition_Updates(t *testing.T) {
	s := &Status{
		Conditions: []Condition{
			NewRenderCondition(ConditionTrue, ReasonRenderSucceeded, ""),
		},
	}

	// Update the existing condition
	updated := NewRenderCondition(ConditionFalse, ReasonRenderFailed, "mutation error")
	s.SetCondition(updated)

	assert.Len(t, s.Conditions, 1)
	assert.Equal(t, ConditionFalse, s.Conditions[0].Status)
	assert.Equal(t, ReasonRenderFailed, s.Conditions[0].Reason)
	assert.Equal(t, "mutation error", s.Conditions[0].Message)
}

func TestStatus_SetCondition_PreservesOtherConditions(t *testing.T) {
	s := &Status{
		Conditions: []Condition{
			{Type: "Other", Status: ConditionTrue, Reason: "SomeReason"},
		},
	}

	cond := NewRenderCondition(ConditionTrue, ReasonRenderSucceeded, "")
	s.SetCondition(cond)

	assert.Len(t, s.Conditions, 2)
	assert.Equal(t, "Other", s.Conditions[0].Type)
	assert.Equal(t, ConditionTypeRendered, s.Conditions[1].Type)
}

// Copyright © 2016 Sidharth Kshatriya
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package engine

import (
	"fmt"
	"strconv"
)

type engineFeatureBool struct {
	Value    bool
	ReadOnly bool
}
type engineFeatureInt struct {
	Value    int
	ReadOnly bool
}
type engineFeatureString struct {
	Value    string
	ReadOnly bool
}

type engineFeatureValue interface {
	Set(value string)
	String() string
}

func (this *engineFeatureBool) Set(value string) {
	if this.ReadOnly {
		panicWith(fmt.Sprintf("Trying assign %v to a read only value: %v", value, this.Value))
	}

	if value == "0" {
		this.Value = false
	} else if value == "1" {
		this.Value = true
	} else {
		panicWith(fmt.Sprintf("Trying to assign a non-boolean value %v to a boolean: %v", value, this.Value))
	}
}

func (this engineFeatureBool) String() string {
	if this.Value {
		return "1"
	}

	return "0"
}

func (this *engineFeatureString) Set(value string) {
	if this.ReadOnly {
		panicWith(fmt.Sprintf("Trying assign %v to a read only value: %v", value, this.Value))
	}
	this.Value = value
}

func (this engineFeatureInt) String() string {
	return strconv.Itoa(this.Value)
}

func (this *engineFeatureInt) Set(value string) {
	if this.ReadOnly {
		panicWith(fmt.Sprintf("Trying assign %v to a read only value: %v", value, this.Value))
	}
	var err error
	this.Value, err = strconv.Atoi(value)
	panicIf(err)
}

func (this engineFeatureString) String() string {
	return this.Value
}

func initFeatureMap() map[string]engineFeatureValue {
	var featureMap = map[string]engineFeatureValue{
		"language_supports_threads": &engineFeatureBool{false, true},
		"language_name":             &engineFeatureString{"PHP", true},
		// @TODO should the exact version be ascertained?
		"language_version":           &engineFeatureString{"7.0", true},
		"encoding":                   &engineFeatureString{"ISO-8859-1", true},
		"protocol_version":           &engineFeatureInt{1, true},
		"supports_async":             &engineFeatureBool{false, true},
		"supports_reverse_debugging": &engineFeatureBool{true, true},
		// @TODO full list
		// "breakpoint_types" : &FeatureString{"line call return exception conditional watch", true},
		"breakpoint_types":    &engineFeatureString{"line", true},
		"multiple_sessions":   &engineFeatureBool{false, false},
		"max_children":        &engineFeatureInt{64, false},
		"max_data":            &engineFeatureInt{2048, false},
		"max_depth":           &engineFeatureInt{1, false},
		"extended_properties": &engineFeatureBool{false, false},
		"show_hidden":         &engineFeatureBool{false, false},
	}

	return featureMap
}

func handleFeatureSet(es *engineState, dCmd dbgpCmd) string {
	n, ok := dCmd.options["n"]
	if !ok {
		panicWith("Please provide -n option in feature_set")
	}

	v, ok := dCmd.options["v"]
	if !ok {
		panicWith("Not provided -v option in feature_set")
	}

	var featureVal engineFeatureValue
	featureVal, ok = es.featureMap[n]
	if !ok {
		panicWith("Unknown option: " + n)
	}

	featureVal.Set(v)
	return fmt.Sprintf(gFeatureSetXMLResponseFormat, dCmd.seqNum, n, 1)
}

func handleFeatureGet(es *engineState, dCmd dbgpCmd) string {
	n, ok := dCmd.options["n"]
	if !ok {
		panicWith("Please provide -n option in feature_get")
	}

	var featureVal engineFeatureValue
	featureVal, ok = es.featureMap[n]
	supported := 1
	if !ok {
		supported = 0
	}

	return fmt.Sprintf(gFeatureGetXMLResponseFormat, dCmd.seqNum, n, supported, featureVal)
}

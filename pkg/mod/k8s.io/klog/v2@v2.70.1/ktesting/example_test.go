/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ktesting_test

import (
	"errors"
	"fmt"
	"time"

	"k8s.io/klog/v2/ktesting"
)

func ExampleUnderlier() {
	logger := ktesting.NewLogger(ktesting.NopTL{}, ktesting.NewConfig(ktesting.Verbosity(4)))

	logger.Error(errors.New("failure"), "I failed", "what", "something")
	logger.WithValues("request", 42).WithValues("anotherValue", "fish").Info("hello world")
	logger.WithValues("request", 42, "anotherValue", "fish").Info("hello world 2", "yetAnotherValue", "thanks")
	logger.WithName("example").Info("with name")
	logger.V(4).Info("higher verbosity")
	logger.V(5).Info("Not captured because of ktesting.Verbosity(4) above. Normally it would be captured because default verbosity is 5.")

	testingLogger, ok := logger.GetSink().(ktesting.Underlier)
	if !ok {
		panic("Should have had a ktesting LogSink!?")
	}

	t := testingLogger.GetUnderlying()
	t.Log("This goes to /dev/null...")

	buffer := testingLogger.GetBuffer()
	fmt.Printf("%s\n", buffer.String())

	log := buffer.Data()
	for i, entry := range log {
		if i > 0 &&
			entry.Timestamp.Sub(log[i-1].Timestamp).Nanoseconds() < 0 {
			fmt.Printf("Unexpected timestamp order: #%d %s > #%d %s", i-1, log[i-1].Timestamp, i, entry.Timestamp)
		}
		// Strip varying time stamp before dumping the struct.
		entry.Timestamp = time.Time{}
		fmt.Printf("log entry #%d: %+v\n", i, entry)
	}

	// Output:
	// ERROR I failed err="failure" what="something"
	// INFO hello world request=42 anotherValue="fish"
	// INFO hello world 2 request=42 anotherValue="fish" yetAnotherValue="thanks"
	// INFO example: with name
	// INFO higher verbosity
	//
	// log entry #0: {Timestamp:0001-01-01 00:00:00 +0000 UTC Type:ERROR Prefix: Message:I failed Verbosity:0 Err:failure WithKVList:[] ParameterKVList:[what something]}
	// log entry #1: {Timestamp:0001-01-01 00:00:00 +0000 UTC Type:INFO Prefix: Message:hello world Verbosity:0 Err:<nil> WithKVList:[request 42 anotherValue fish] ParameterKVList:[]}
	// log entry #2: {Timestamp:0001-01-01 00:00:00 +0000 UTC Type:INFO Prefix: Message:hello world 2 Verbosity:0 Err:<nil> WithKVList:[request 42 anotherValue fish] ParameterKVList:[yetAnotherValue thanks]}
	// log entry #3: {Timestamp:0001-01-01 00:00:00 +0000 UTC Type:INFO Prefix:example Message:with name Verbosity:0 Err:<nil> WithKVList:[] ParameterKVList:[]}
	// log entry #4: {Timestamp:0001-01-01 00:00:00 +0000 UTC Type:INFO Prefix: Message:higher verbosity Verbosity:4 Err:<nil> WithKVList:[] ParameterKVList:[]}
}

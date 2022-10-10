/*
Copyright 2019 The Kubernetes Authors.
Copyright 2020 Intel Coporation.

SPDX-License-Identifier: Apache-2.0
*/

package ktesting_test

import (
	"context"
	"testing"

	"k8s.io/klog/v2"
	"k8s.io/klog/v2/ktesting"
)

func TestContextual(t *testing.T) {
	logger, ctx := ktesting.NewTestContext(t)

	doSomething(ctx)

	// When contextual logging is disabled, the output goes to klog
	// instead of the testing logger.
	state := klog.CaptureState()
	defer state.Restore()
	klog.EnableContextualLogging(false)
	doSomething(ctx)

	testingLogger, ok := logger.GetSink().(ktesting.Underlier)
	if !ok {
		t.Fatal("Should have had a ktesting LogSink!?")
	}

	actual := testingLogger.GetBuffer().String()
	expected := `INFO hello world
INFO foo: hello also from me
`
	if actual != expected {
		t.Errorf("mismatch in captured output, expected:\n%s\ngot:\n%s\n", expected, actual)
	}
}

func doSomething(ctx context.Context) {
	logger := klog.FromContext(ctx)
	logger.Info("hello world")

	logger = logger.WithName("foo")
	ctx = klog.NewContext(ctx, logger)
	doSomeMore(ctx)
}

func doSomeMore(ctx context.Context) {
	logger := klog.FromContext(ctx)
	logger.Info("hello also from me")
}

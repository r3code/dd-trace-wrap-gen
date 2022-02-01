package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func TestTestInterfaceWithOpenTelemetryTracing_F(t *testing.T) {
	t.Run("it creates span without error mark", func(t *testing.T) {
		mt := mocktracer.Start()
		defer mt.Stop()

		impl := &testImpl{result1: "1", result2: "OK"}
		noopSpanDecorator := func(span tracer.Span, params, results map[string]interface{}) {
		}
		wrapped := NewTestInterfaceWithTracing(impl, "test.operation", noopSpanDecorator)

		span, ctx := tracer.StartSpanFromContext(context.Background(), "root_span")
		r1, r2, err := wrapped.F(ctx, "a1", "a2")
		span.Finish()
		require.NoError(t, err)
		require.Equal(t, "1", r1)
		require.Equal(t, "OK", r2)
		spans := mt.FinishedSpans()
		require.Equal(t, 2, len(spans))
		require.Equal(t, "test.operation", spans[0].OperationName())
		require.Equal(t, "root_span", spans[1].OperationName())
	})

	t.Run("it marks span with error", func(t *testing.T) {
		mt := mocktracer.Start()
		defer mt.Stop()

		impl := &testImpl{err: errors.New("it is a test error")}
		noopSpanDecorator := func(span tracer.Span, params, results map[string]interface{}) {
		}
		wrapped := NewTestInterfaceWithTracing(impl, "test.operation", noopSpanDecorator)

		_, ctx := tracer.StartSpanFromContext(context.Background(), "root_span")
		_, _, err := wrapped.F(ctx, "a1", "a2")
		require.Error(t, err)
		spans := mt.FinishedSpans()
		require.Equal(t, 1, len(spans))
		require.Equal(t, "test.operation", spans[0].OperationName())
		require.Contains(t, fmt.Sprintf("%v", spans[0].Tag(ext.Error)), impl.err.Error())
	})

	t.Run("it decorates span with tags", func(t *testing.T) {
		mt := mocktracer.Start()
		defer mt.Stop()

		impl := &testImpl{result1: "1", result2: "OK"}
		spanDecorator := func(span tracer.Span, params, results map[string]interface{}) {
			for name, val := range results {
				span.SetTag(name, val)
			}
		}
		wrapped := NewTestInterfaceWithTracing(impl, "test.operation", spanDecorator)

		_, ctx := tracer.StartSpanFromContext(context.Background(), "root_span")
		_, _, err := wrapped.F(ctx, "a1", "a2")
		require.NoError(t, err)
		spans := mt.FinishedSpans()
		require.Equal(t, 1, len(spans))
		require.Equal(t, "test.operation", spans[0].OperationName())
		require.Equal(t, impl.result1, fmt.Sprintf("%v", spans[0].Tag("result1")))
		require.Equal(t, impl.result2, fmt.Sprintf("%v", spans[0].Tag("result2")))
	})

	t.Run("it decides to skip marking span with error", func(t *testing.T) {
		mt := mocktracer.Start()
		defer mt.Stop()

		impl := &testImpl{err: errors.New("skip me please")}
		noopSpanDecorator := func(span tracer.Span, params, results map[string]interface{}) {}
		checkErrorMarkRequired := func(err error) bool {
			if err.Error() == "skip me please" {
				return false
			}

			if err.Error() == "skip me please, case 2" {
				return false
			}

			return true
		}
		wrapped := NewTestInterfaceWithTracing(impl, "test.operation", noopSpanDecorator, checkErrorMarkRequired)

		_, ctx := tracer.StartSpanFromContext(context.Background(), "root_span")
		_, _, err := wrapped.F(ctx, "a1", "a2")
		require.Error(t, err)
		spans := mt.FinishedSpans()
		require.Equal(t, 1, len(spans))
		require.Equal(t, "test.operation", spans[0].OperationName())
		require.NotContains(t, fmt.Sprintf("%v", spans[0].Tag(ext.Error)), impl.err.Error())
	})
}

type testImpl struct {
	result1, result2 string
	err              error
}

func (t testImpl) F(_ context.Context, _ string, _ ...string) (result1, result2 string, err error) {
	if t.err != nil {
		return "", "", t.err
	}
	return t.result1, t.result2, nil
}

func (t testImpl) NoError(s string) string {
	return s
}

func (t testImpl) NoParamsOrResults() {}

func (t testImpl) Channels(_ chan bool, _ chan<- bool, _ <-chan bool) {}

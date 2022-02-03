package tests

// Code generated by github.com/r3code/dd-trace-wrap-gen tool. DO NOT EDIT!
//

import (
	"context"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// TestInterfaceWithTracing implements TestInterface interface instrumented with datadog spans
type TestInterfaceWithTracing struct {
	TestInterface
	_spanName         string
	_spanDecorator    func(span tracer.Span, params, results map[string]interface{})
	_errorMarkDecider func(err error) bool
}

// NewTestInterfaceWithTracing returns TestInterfaceWithTracing for the base service with a specified
// Datadog’s span name spanName (equals OpenTracing “component” tag), and allows to add extra data to the span by spanDecorator (a func to add some extra tags for a span).
// Pass nil if you don't need decorations.
// You can skip marking a span with an error mark by returning false in errorMarkDecider func. Optional, by default the decider always returns true.
//
//  Note: when using Datadog, the OpenTracing operation name is a resource and the OpenTracing “component” tag is Datadog’s span name.
//  SpanName in DataDog becomes an "operation name" and "resource name" is taken from $method.Name
//  Example. Create a span for a http request for url /user/profile:
//    spanName = "http.request"
//    resource = "/user/profile"
func NewTestInterfaceWithTracing(base TestInterface, spanName string, spanDecorator func(span tracer.Span, params, results map[string]interface{}), errorMarkDecider ...func(err error) bool) TestInterfaceWithTracing {
	d := TestInterfaceWithTracing{
		TestInterface:     base,
		_spanName:         spanName,
		_errorMarkDecider: func(err error) bool { return true }, // by default always allow mark a span having an error
	}
	if spanDecorator != nil {
		d._spanDecorator = spanDecorator
	}

	if len(errorMarkDecider) > 0 && errorMarkDecider[0] != nil {
		d._errorMarkDecider = errorMarkDecider[0]
	}

	return d
}

// F implements TestInterface
func (_d TestInterfaceWithTracing) F(ctx context.Context, a1 string, a2 ...string) (result1 string, result2 string, err error) {
	_span, ctx := tracer.StartSpanFromContext(ctx, _d._spanName, tracer.ResourceName("F"))
	defer func() {
		if _d._spanDecorator != nil {
			_d._spanDecorator(_span, map[string]interface{}{
				"ctx": ctx,
				"a1":  a1,
				"a2":  a2}, map[string]interface{}{
				"result1": result1,
				"result2": result2,
				"err":     err})
		}
		var opts []tracer.FinishOption
		if err != nil && _d._errorMarkDecider(err) {
			opts = append(opts, tracer.WithError(err))
		}
		_span.Finish(opts...)
	}()
	return _d.TestInterface.F(ctx, a1, a2...)
}

// NoErrorWithContext implements TestInterface
func (_d TestInterfaceWithTracing) NoErrorWithContext(ctx context.Context, s1 string) (s2 string) {
	_span, ctx := tracer.StartSpanFromContext(ctx, _d._spanName, tracer.ResourceName("NoErrorWithContext"))
	defer func() {
		if _d._spanDecorator != nil {
			_d._spanDecorator(_span, map[string]interface{}{
				"ctx": ctx,
				"s1":  s1}, map[string]interface{}{
				"s2": s2})
		}
		_span.Finish()

	}()
	return _d.TestInterface.NoErrorWithContext(ctx, s1)
}

// Package gateway provides a drop-in replacement for net/http.ListenAndServe for use in AWS Lambda & API Gateway.
package gateway

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/getas/alb-gateway/events"
)

// ListenAndServe is a drop-in replacement for
// http.ListenAndServe for use within AWS Lambda.
//
// ListenAndServe always returns a non-nil error.
func ListenAndServe(addr string, h http.Handler) error {
	if h == nil {
		h = http.DefaultServeMux
	}

	lambda.Start(func(ctx context.Context, e events.LambdaTargetGroupRequest) (events.LambdaTargetGroupResponse, error) {
		r, err := NewRequest(ctx, e)
		if err != nil {
			return events.LambdaTargetGroupResponse{}, err
		}

		w := NewResponse()
		h.ServeHTTP(w, r)
		return w.End(), nil
	})

	return nil
}

package gmongo

import (
	"context"
	"testing"
)

func TestCollectionQueryWithContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), struct{}{}, "trace")
	collection := &Collection{}

	query := collection.QueryWithContext(ctx)

	if query.ctx != ctx {
		t.Fatalf("expected query builder context to be set")
	}
}

func TestNewAggregationWithContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), struct{}{}, "trace")
	collection := &Collection{}

	aggregation := NewAggregationWithContext(ctx, collection)

	if aggregation.ctx != ctx {
		t.Fatalf("expected aggregation builder context to be set")
	}
}

func TestCollectionNewAggregationWithContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), struct{}{}, "trace")
	collection := &Collection{}

	aggregation := collection.NewAggregationWithContext(ctx)

	if aggregation.ctx != ctx {
		t.Fatalf("expected aggregation builder context to be set")
	}
}

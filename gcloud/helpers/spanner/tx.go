package spanner

import (
	"context"

	"cloud.google.com/go/spanner"
	"cloud.google.com/go/spanner/apiv1/spannerpb"
)

// Tx is a common interface for spanner read only and read/write transactions that only read
type Tx interface {
	Read(ctx context.Context, table string, keys spanner.KeySet, columns []string) *spanner.RowIterator
	ReadUsingIndex(ctx context.Context, table, index string, keys spanner.KeySet, columns []string) (ri *spanner.RowIterator)
	ReadWithOptions(ctx context.Context, table string, keys spanner.KeySet, columns []string, opts *spanner.ReadOptions) (ri *spanner.RowIterator)
	ReadRow(ctx context.Context, table string, key spanner.Key, columns []string) (*spanner.Row, error)
	Query(ctx context.Context, statement spanner.Statement) *spanner.RowIterator
	QueryWithStats(ctx context.Context, statement spanner.Statement) *spanner.RowIterator
	AnalyzeQuery(ctx context.Context, statement spanner.Statement) (*spannerpb.QueryPlan, error)
	BufferWrite([]*spanner.Mutation) error
	Update(ctx context.Context, stmt spanner.Statement) (rowCount int64, err error)
}

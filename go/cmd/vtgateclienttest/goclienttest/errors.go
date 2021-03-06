// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goclienttest

import (
	"io"
	"testing"

	"golang.org/x/net/context"

	"github.com/youtube/vitess/go/sqltypes"

	"github.com/youtube/vitess/go/vt/vterrors"
	"github.com/youtube/vitess/go/vt/vtgate/vtgateconn"

	querypb "github.com/youtube/vitess/go/vt/proto/query"
	vtgatepb "github.com/youtube/vitess/go/vt/proto/vtgate"
	vtrpcpb "github.com/youtube/vitess/go/vt/proto/vtrpc"
)

var (
	errorPrefix        = "error://"
	partialErrorPrefix = "partialerror://"

	executeErrors = map[string]vtrpcpb.ErrorCode{
		"bad input":         vtrpcpb.ErrorCode_BAD_INPUT,
		"deadline exceeded": vtrpcpb.ErrorCode_DEADLINE_EXCEEDED,
		"integrity error":   vtrpcpb.ErrorCode_INTEGRITY_ERROR,
		"transient error":   vtrpcpb.ErrorCode_TRANSIENT_ERROR,
		"unauthenticated":   vtrpcpb.ErrorCode_UNAUTHENTICATED,
		"aborted":           vtrpcpb.ErrorCode_NOT_IN_TX,
		"unknown error":     vtrpcpb.ErrorCode_UNKNOWN_ERROR,
	}
)

// testErrors exercises the test cases provided by the "errors" service.
func testErrors(t *testing.T, conn *vtgateconn.VTGateConn) {
	testExecuteErrors(t, conn)
	testStreamExecuteErrors(t, conn)
	testTransactionExecuteErrors(t, conn)
	testUpdateStreamErrors(t, conn)
}

func testExecuteErrors(t *testing.T, conn *vtgateconn.VTGateConn) {
	ctx := context.Background()

	checkExecuteErrors(t, func(query string) error {
		_, err := conn.Execute(ctx, query, bindVars, tabletType)
		return err
	})
	checkExecuteErrors(t, func(query string) error {
		_, err := conn.ExecuteShards(ctx, query, keyspace, shards, bindVars, tabletType)
		return err
	})
	checkExecuteErrors(t, func(query string) error {
		_, err := conn.ExecuteKeyspaceIds(ctx, query, keyspace, keyspaceIDs, bindVars, tabletType)
		return err
	})
	checkExecuteErrors(t, func(query string) error {
		_, err := conn.ExecuteKeyRanges(ctx, query, keyspace, keyRanges, bindVars, tabletType)
		return err
	})
	checkExecuteErrors(t, func(query string) error {
		_, err := conn.ExecuteEntityIds(ctx, query, keyspace, "column1", entityKeyspaceIDs, bindVars, tabletType)
		return err
	})
	checkExecuteErrors(t, func(query string) error {
		_, err := conn.ExecuteBatchShards(ctx, []*vtgatepb.BoundShardQuery{
			{
				Query: &querypb.BoundQuery{
					Sql:           query,
					BindVariables: bindVarsP3,
				},
				Keyspace: keyspace,
				Shards:   shards,
			},
		}, tabletType, true)
		return err
	})
	checkExecuteErrors(t, func(query string) error {
		_, err := conn.ExecuteBatchKeyspaceIds(ctx, []*vtgatepb.BoundKeyspaceIdQuery{
			{
				Query: &querypb.BoundQuery{
					Sql:           query,
					BindVariables: bindVarsP3,
				},
				Keyspace:    keyspace,
				KeyspaceIds: keyspaceIDs,
			},
		}, tabletType, true)
		return err
	})
}

func testStreamExecuteErrors(t *testing.T, conn *vtgateconn.VTGateConn) {
	ctx := context.Background()

	checkStreamExecuteErrors(t, func(query string) error {
		return getStreamError(conn.StreamExecute(ctx, query, bindVars, tabletType))
	})
	checkStreamExecuteErrors(t, func(query string) error {
		return getStreamError(conn.StreamExecuteShards(ctx, query, keyspace, shards, bindVars, tabletType))
	})
	checkStreamExecuteErrors(t, func(query string) error {
		return getStreamError(conn.StreamExecuteKeyspaceIds(ctx, query, keyspace, keyspaceIDs, bindVars, tabletType))
	})
	checkStreamExecuteErrors(t, func(query string) error {
		return getStreamError(conn.StreamExecuteKeyRanges(ctx, query, keyspace, keyRanges, bindVars, tabletType))
	})
}

func testUpdateStreamErrors(t *testing.T, conn *vtgateconn.VTGateConn) {
	ctx := context.Background()

	checkStreamExecuteErrors(t, func(query string) error {
		return getUpdateStreamError(conn.UpdateStream(ctx, query, nil, tabletType, 0, nil))
	})
}

func testTransactionExecuteErrors(t *testing.T, conn *vtgateconn.VTGateConn) {
	ctx := context.Background()

	checkTransactionExecuteErrors(t, conn, func(tx *vtgateconn.VTGateTx, query string) error {
		_, err := tx.Execute(ctx, query, bindVars, tabletType)
		return err
	})
	checkTransactionExecuteErrors(t, conn, func(tx *vtgateconn.VTGateTx, query string) error {
		_, err := tx.ExecuteShards(ctx, query, keyspace, shards, bindVars, tabletType)
		return err
	})
	checkTransactionExecuteErrors(t, conn, func(tx *vtgateconn.VTGateTx, query string) error {
		_, err := tx.ExecuteKeyspaceIds(ctx, query, keyspace, keyspaceIDs, bindVars, tabletType)
		return err
	})
	checkTransactionExecuteErrors(t, conn, func(tx *vtgateconn.VTGateTx, query string) error {
		_, err := tx.ExecuteKeyRanges(ctx, query, keyspace, keyRanges, bindVars, tabletType)
		return err
	})
	checkTransactionExecuteErrors(t, conn, func(tx *vtgateconn.VTGateTx, query string) error {
		_, err := tx.ExecuteEntityIds(ctx, query, keyspace, "column1", entityKeyspaceIDs, bindVars, tabletType)
		return err
	})
	checkTransactionExecuteErrors(t, conn, func(tx *vtgateconn.VTGateTx, query string) error {
		_, err := tx.ExecuteBatchShards(ctx, []*vtgatepb.BoundShardQuery{
			{
				Query: &querypb.BoundQuery{
					Sql:           query,
					BindVariables: bindVarsP3,
				},
				Keyspace: keyspace,
				Shards:   shards,
			},
		}, tabletType)
		return err
	})
	checkTransactionExecuteErrors(t, conn, func(tx *vtgateconn.VTGateTx, query string) error {
		_, err := tx.ExecuteBatchKeyspaceIds(ctx, []*vtgatepb.BoundKeyspaceIdQuery{
			{
				Query: &querypb.BoundQuery{
					Sql:           query,
					BindVariables: bindVarsP3,
				},
				Keyspace:    keyspace,
				KeyspaceIds: keyspaceIDs,
			},
		}, tabletType)
		return err
	})
}

func getStreamError(stream sqltypes.ResultStream, err error) error {
	if err != nil {
		return err
	}
	for {
		_, err := stream.Recv()
		switch err {
		case nil:
			// keep going
		case io.EOF:
			return nil
		default:
			return err
		}
	}
}

func getUpdateStreamError(stream vtgateconn.UpdateStreamReader, err error) error {
	if err != nil {
		return err
	}
	for {
		_, _, err := stream.Recv()
		switch err {
		case nil:
			// keep going
		case io.EOF:
			return nil
		default:
			return err
		}
	}
}

func checkExecuteErrors(t *testing.T, execute func(string) error) {
	for errStr, errCode := range executeErrors {
		query := errorPrefix + errStr
		checkError(t, execute(query), query, errStr, errCode)

		query = partialErrorPrefix + errStr
		checkError(t, execute(query), query, errStr, errCode)
	}
}

func checkStreamExecuteErrors(t *testing.T, execute func(string) error) {
	for errStr, errCode := range executeErrors {
		query := errorPrefix + errStr
		checkError(t, execute(query), query, errStr, errCode)
	}
}

func checkTransactionExecuteErrors(t *testing.T, conn *vtgateconn.VTGateConn, execute func(tx *vtgateconn.VTGateTx, query string) error) {
	ctx := context.Background()

	for errStr, errCode := range executeErrors {
		query := errorPrefix + errStr
		tx, err := conn.Begin(ctx)
		if err != nil {
			t.Errorf("[%v] Begin error: %v", query, err)
		}
		checkError(t, execute(tx, query), query, errStr, errCode)

		// Partial error where server doesn't close the session.
		query = partialErrorPrefix + errStr
		tx, err = conn.Begin(ctx)
		if err != nil {
			t.Errorf("[%v] Begin error: %v", query, err)
		}
		checkError(t, execute(tx, query), query, errStr, errCode)
		// The transaction should still be usable now.
		if err := tx.Rollback(ctx); err != nil {
			t.Errorf("[%v] Rollback error: %v", query, err)
		}

		// Partial error where server closes the session.
		tx, err = conn.Begin(ctx)
		if err != nil {
			t.Errorf("[%v] Begin error: %v", query, err)
		}
		query = partialErrorPrefix + errStr + "/close transaction"
		checkError(t, execute(tx, query), query, errStr, errCode)
		// The transaction should be unusable now.
		if tx.Rollback(ctx) == nil {
			t.Errorf("[%v] expected Rollback error, got nil", query)
		}
	}
}

func checkError(t *testing.T, err error, query, errStr string, errCode vtrpcpb.ErrorCode) {
	if err == nil {
		t.Errorf("[%v] expected error, got nil", query)
		return
	}
	switch vtErr := err.(type) {
	case *vterrors.VitessError:
		if got, want := vtErr.VtErrorCode(), errCode; got != want {
			t.Errorf("[%v] error code = %v, want %v", query, got, want)
		}
	default:
		t.Errorf("[%v] unrecognized error type: %T, error: %#v", query, err, err)
		return
	}

}

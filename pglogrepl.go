package pglogrepl

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgconn"
	errors "golang.org/x/xerrors"
)

type IdentifySystemResult struct {
	SystemID string
	Timeline int32
	XlogPos  string
	DBName   string
}

func IdentifySystem(ctx context.Context, conn *pgconn.PgConn) (IdentifySystemResult, error) {
	return ParseIdentifySystem(conn.Exec(ctx, "IDENTIFY_SYSTEM"))
}

// ParseIdentifySystem parses the result of the IDENTIFY_SYSTEM command.
func ParseIdentifySystem(mrr *pgconn.MultiResultReader) (IdentifySystemResult, error) {
	var isr IdentifySystemResult
	results, err := mrr.ReadAll()
	if err != nil {
		return isr, err
	}

	if len(results) != 1 {
		return isr, errors.Errorf("expected 1 result set, got %d", len(results))
	}

	result := results[0]
	if len(result.Rows) != 1 {
		return isr, errors.Errorf("expected 1 result row, got %d", len(result.Rows))
	}

	row := result.Rows[0]
	if len(row) != 4 {
		return isr, errors.Errorf("expected 4 result columns, got %d", len(row))
	}

	isr.SystemID = string(row[0])
	timeline, err := strconv.ParseInt(string(row[1]), 10, 32)
	if err != nil {
		return isr, errors.Errorf("failed to parse timeline: %w", err)
	}
	isr.Timeline = int32(timeline)
	isr.XlogPos = string(row[2])
	isr.DBName = string(row[3])

	return isr, nil
}

type CreateReplicationSlotOptions struct {
	Temporary      bool
	SnapshotAction string
}

type CreateReplicationSlotResult struct {
	SlotName        string
	ConsistentPoint string
	SnapshotName    string
	OutputPlugin    string
}

// CreateReplicationSlot creates a logical replication slot.
func CreateReplicationSlot(
	ctx context.Context,
	conn *pgconn.PgConn,
	slotName string,
	outputPlugin string,
	options CreateReplicationSlotOptions,
) (CreateReplicationSlotResult, error) {
	temporaryString := "TEMPORARY"
	sql := fmt.Sprintf("CREATE_REPLICATION_SLOT %s %s LOGICAL %s %s", slotName, temporaryString, outputPlugin, options.SnapshotAction)
	return ParseCreateReplicationSlot(conn.Exec(ctx, sql))
}

// ParseCreateReplicationSlot parses the result of the CREATE_REPLICATION_SLOT command.
func ParseCreateReplicationSlot(mrr *pgconn.MultiResultReader) (CreateReplicationSlotResult, error) {
	var crsr CreateReplicationSlotResult
	results, err := mrr.ReadAll()
	if err != nil {
		return crsr, err
	}

	if len(results) != 1 {
		return crsr, errors.Errorf("expected 1 result set, got %d", len(results))
	}

	result := results[0]
	if len(result.Rows) != 1 {
		return crsr, errors.Errorf("expected 1 result row, got %d", len(result.Rows))
	}

	row := result.Rows[0]
	if len(row) != 4 {
		return crsr, errors.Errorf("expected 4 result columns, got %d", len(row))
	}

	crsr.SlotName = string(row[0])
	crsr.ConsistentPoint = string(row[1])
	crsr.SnapshotName = string(row[2])
	crsr.OutputPlugin = string(row[3])

	return crsr, nil
}

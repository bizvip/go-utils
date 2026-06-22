package snowflake_test

import (
	"testing"

	"github.com/bizvip/go-utils/base/snowflake"
)

func TestShortIdGeneratorWorkerIDRange(t *testing.T) {
	gen, err := snowflake.NewShortIdGenerator(31)
	if err != nil {
		t.Fatalf("NewShortIdGenerator(31) failed: %v", err)
	}

	id, err := gen.NextID()
	if err != nil {
		t.Fatalf("NextID failed: %v", err)
	}

	_, workerID, _ := gen.Decompose(id)
	if workerID != 31 {
		t.Fatalf("decomposed workerID = %d, want 31", workerID)
	}

	if _, err := snowflake.NewShortIdGenerator(32); err == nil {
		t.Fatalf("expected error for workerId out of range (max 31)")
	}
}

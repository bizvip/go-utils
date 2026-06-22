package bizid

import (
	"errors"
	"strings"
	"sync"
	"testing"
	"time"
)

func resetForTest() {
	mu.Lock()
	gen = nil
	mu.Unlock()
}

func TestEncodeBase26(t *testing.T) {
	cases := []struct {
		in   uint64
		want string
	}{
		{0, "A"},
		{1, "B"},
		{25, "Z"},
		{26, "BA"},
		{27, "BB"},
		{676, "BAA"}, // 26*26
	}
	for _, tc := range cases {
		if got := encodeBase26(tc.in); got != tc.want {
			t.Fatalf("encodeBase26(%d) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestTimeBase26IsStableWithinSameMinute(t *testing.T) {
	tt := time.Date(2026, 5, 22, 14, 30, 0, 0, time.UTC)
	first := timeBase26(tt)
	other := timeBase26(tt.Add(40 * time.Second))
	if first != other {
		t.Fatalf("same minute should yield equal base26 segment, got %q vs %q", first, other)
	}
	if first == timeBase26(tt.Add(time.Minute)) {
		t.Fatalf("different minute should yield different base26 segment")
	}
}

func TestGetSnowflakeIDBeforeInitReturnsError(t *testing.T) {
	resetForTest()
	_, err := GetSnowflakeID()
	if !errors.Is(err, ErrNotInitialized) {
		t.Fatalf("expected ErrNotInitialized, got %v", err)
	}
}

func TestMustGetSnowflakeIDBeforeInitPanics(t *testing.T) {
	resetForTest()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic when MustGetSnowflakeID called before Init")
		}
	}()
	_ = MustGetSnowflakeID()
}

func TestNewBeforeInitPanics(t *testing.T) {
	resetForTest()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic when New called before Init")
		}
	}()
	_ = New("RG")
}

func TestInitRejectsBadWorkerID(t *testing.T) {
	resetForTest()
	if err := Init(31); err != nil {
		t.Fatalf("expected max workerId to be accepted: %v", err)
	}
	if err := Init(32); err == nil {
		t.Fatalf("expected error for workerId out of range (max 31)")
	}
}

func TestInitAndGetSnowflakeID(t *testing.T) {
	resetForTest()
	if err := Init(0); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	id, err := GetSnowflakeID()
	if err != nil {
		t.Fatalf("GetSnowflakeID failed: %v", err)
	}
	if id <= 0 {
		t.Fatalf("expected positive id, got %d", id)
	}
}

func TestNewIncludesPrefixAndSnowflake(t *testing.T) {
	resetForTest()
	if err := Init(0); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	id := New("RG")
	if !strings.HasPrefix(id, "RG") {
		t.Fatalf("expected prefix RG, got %q", id)
	}
	if len(id) < 4 {
		t.Fatalf("unexpectedly short id: %q", id)
	}
	last := id[len(id)-1]
	if last < '0' || last > '9' {
		t.Fatalf("expected id to end with digit, got %q", id)
	}
}

func TestShortcutsUseCorrectPrefix(t *testing.T) {
	resetForTest()
	if err := Init(0); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	cases := []struct {
		prefix string
		fn     func() string
	}{
		{"RG", GetRechargeOrderID},
		{"WD", GetWithdrawID},
		{"DP", GetDepositOrderID},
		{"BU", GetBuyOrderID},
		{"EX", GetExchangeID},
		{"BO", GetBonusOrderID},
		{"GM", GetGemOrderID},
		{"GD", GetGoldOrderID},
	}
	for _, tc := range cases {
		id := tc.fn()
		if !strings.HasPrefix(id, tc.prefix) {
			t.Fatalf("shortcut for %q produced %q, prefix mismatch", tc.prefix, id)
		}
	}
}

func TestNewIDsAreDistinctConcurrent(t *testing.T) {
	resetForTest()
	if err := Init(0); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	const (
		workers = 8
		perWk   = 200
	)
	var (
		wg   sync.WaitGroup
		lock sync.Mutex
		seen = make(map[string]struct{}, workers*perWk)
	)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			local := make([]string, 0, perWk)
			for i := 0; i < perWk; i++ {
				local = append(local, New("RG"))
			}
			lock.Lock()
			for _, id := range local {
				if _, dup := seen[id]; dup {
					t.Errorf("duplicate id: %s", id)
				}
				seen[id] = struct{}{}
			}
			lock.Unlock()
		}()
	}
	wg.Wait()
}

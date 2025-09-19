package snowflake

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const (
	customEpoch int64 = 1758240000000 // 2025-09-19

	// Bit allocation
	timeBits     uint = 39 // ~17.4 years
	workerIdBits uint = 4  // up to 16 nodes (0-15)
	sequenceBits uint = 4  // ~ 16000/s
)

const spinSleep = 50 * time.Microsecond

// Precomputed masks and shifts (all in uint64 for safer bit ops)
var (
	maxWorkerId uint64 = (1 << workerIdBits) - 1
	maxSequence uint64 = (1 << sequenceBits) - 1

	workerIdShift  uint = sequenceBits
	timestampShift uint = sequenceBits + workerIdBits
	totalBits           = timeBits + workerIdBits + sequenceBits
)

// ShortIdGenerator is a monotonic, short-length Snowflake-like ID generator.
// It uses a process-local monotonic clock plus a startup offset derived from wall clock.
type ShortIdGenerator struct {
	mu sync.Mutex

	// Last emitted time (milliseconds since customEpoch, monotonic within the process)
	lastTime uint64

	// Configuration
	workerId uint64

	// Per-millisecond rolling sequence
	sequence uint64

	// Monotonic time base (process start moment)
	baseMono time.Time
	// Startup offset: (wall clock in ms since epoch) - customEpoch
	baseOffsetMs uint64
}

// NewShortIdGenerator creates a new generator for the given workerId.
func NewShortIdGenerator(workerId int64) (*ShortIdGenerator, error) {
	// Defensive: ensure total width fits in uint64
	//if totalBits > 63 {
	//	return nil, fmt.Errorf("invalid bit allocation: time=%d worker=%d seq=%d (sum=%d > 63)",
	//		timeBits, workerIdBits, sequenceBits, totalBits)
	//}

	// Validate workerId within bit-derived range
	if workerId < 0 || uint64(workerId) > maxWorkerId {
		return nil, fmt.Errorf("worker ID must be between 0 and %d", maxWorkerId)
	}

	// Derive startup offset using wall clock (can be <0 if before epoch; clamp to 0)
	now := time.Now()
	baseOffset := now.UnixMilli() - customEpoch
	var baseOffsetMs uint64
	if baseOffset > 0 {
		baseOffsetMs = uint64(baseOffset)
	} else {
		// If wall clock is before custom epoch, start from 0 but still monotonic within process
		baseOffsetMs = 0
	}

	return &ShortIdGenerator{
		workerId:     uint64(workerId),
		baseMono:     now,
		baseOffsetMs: baseOffsetMs,
	}, nil
}

// monoNowMs returns a process-monotonic "milliseconds since customEpoch".
func (g *ShortIdGenerator) monoNowMs() uint64 {
	elapsed := time.Since(g.baseMono).Milliseconds()
	if elapsed < 0 {
		// Should not happen; keep it defensive.
		return g.baseOffsetMs
	}
	return g.baseOffsetMs + uint64(elapsed)
}

// NextID generates one ID. It is strictly monotonic within the process for the same worker.
func (g *ShortIdGenerator) NextID() (uint64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := g.monoNowMs()

	// Monotonic guarantee within the process
	if now < g.lastTime {
		// Should not happen with monotonic time; clamp to lastTime
		now = g.lastTime
	}

	if now == g.lastTime {
		// Same millisecond: advance sequence
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			// Sequence exhausted: wait for next millisecond
			for {
				runtime.Gosched()
				time.Sleep(spinSleep)
				now = g.monoNowMs()
				if now > g.lastTime {
					break
				}
			}
		}
	} else {
		// New millisecond: reset sequence
		g.sequence = 0
	}

	g.lastTime = now

	id := (now << timestampShift) | (g.workerId << workerIdShift) | g.sequence
	return id, nil
}

// BatchNext generates n IDs in one critical section to reduce lock contention.
func (g *ShortIdGenerator) BatchNext(n int) ([]uint64, error) {
	if n <= 0 {
		return nil, nil
	}
	out := make([]uint64, n)

	g.mu.Lock()
	defer g.mu.Unlock()

	for i := 0; i < n; i++ {
		now := g.monoNowMs()
		if now < g.lastTime {
			now = g.lastTime
		}

		if now == g.lastTime {
			g.sequence = (g.sequence + 1) & maxSequence
			if g.sequence == 0 {
				for {
					runtime.Gosched()
					time.Sleep(spinSleep)
					now = g.monoNowMs()
					if now > g.lastTime {
						break
					}
				}
			}
		} else {
			g.sequence = 0
		}

		g.lastTime = now
		out[i] = (now << timestampShift) | (g.workerId << workerIdShift) | g.sequence
	}
	return out, nil
}

// Decompose splits an ID into (timestampMsSinceCustomEpoch, workerId, sequence).
func (g *ShortIdGenerator) Decompose(id uint64) (tsMs uint64, workerId uint64, seq uint64) {
	seq = id & maxSequence
	workerId = (id >> workerIdShift) & maxWorkerId
	tsMs = id >> timestampShift
	return
}

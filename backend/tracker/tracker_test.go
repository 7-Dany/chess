package tracker

import (
	"testing"
)

// TestHashPositionTracker verifies that the tracker correctly counts position
// occurrences and supports undo — the foundation of threefold repetition
// detection.
func TestHashPositionTracker(t *testing.T) {
	// =========================================================================
	// Basic counting — Record and Count.
	// =========================================================================

	t.Run("a fresh tracker reports count 0 for any hash", func(t *testing.T) {
		tr := NewPositionTracker()
		if got := tr.Count(42); got != 0 {
			t.Errorf("Count(42) = %d, want 0", got)
		}
		if got := tr.Count(0); got != 0 {
			t.Errorf("Count(0) = %d, want 0", got)
		}
	})

	t.Run("recording a hash once gives a count of 1", func(t *testing.T) {
		tr := NewPositionTracker()
		tr.Record(42)
		if got := tr.Count(42); got != 1 {
			t.Errorf("Count(42) = %d, want 1", got)
		}
	})

	t.Run("recording the same hash three times gives a count of 3 (threefold)", func(t *testing.T) {
		tr := NewPositionTracker()
		tr.Record(42)
		tr.Record(42)
		tr.Record(42)
		if got := tr.Count(42); got != 3 {
			t.Errorf("Count(42) = %d, want 3", got)
		}
	})

	t.Run("different hashes have independent counts", func(t *testing.T) {
		tr := NewPositionTracker()
		tr.Record(1)
		tr.Record(2)
		tr.Record(2)
		tr.Record(3)
		tr.Record(3)
		tr.Record(3)
		if got := tr.Count(1); got != 1 {
			t.Errorf("Count(1) = %d, want 1", got)
		}
		if got := tr.Count(2); got != 2 {
			t.Errorf("Count(2) = %d, want 2", got)
		}
		if got := tr.Count(3); got != 3 {
			t.Errorf("Count(3) = %d, want 3", got)
		}
	})

	// =========================================================================
	// Undo — decrementing counts.
	// =========================================================================

	t.Run("undo decrements the count by 1", func(t *testing.T) {
		tr := NewPositionTracker()
		tr.Record(42)
		tr.Record(42)
		tr.Record(42)
		tr.Undo(42)
		if got := tr.Count(42); got != 2 {
			t.Errorf("Count(42) after undo = %d, want 2", got)
		}
	})

	t.Run("undo back to zero removes the entry (count is 0)", func(t *testing.T) {
		tr := NewPositionTracker()
		tr.Record(42)
		tr.Undo(42)
		if got := tr.Count(42); got != 0 {
			t.Errorf("Count(42) after undo to 0 = %d, want 0", got)
		}
	})

	t.Run("undo a hash that was never recorded is a no-op (count stays 0)", func(t *testing.T) {
		tr := NewPositionTracker()
		tr.Undo(42) // should not panic
		if got := tr.Count(42); got != 0 {
			t.Errorf("Count(42) after undo of unrecorded hash = %d, want 0", got)
		}
	})

	t.Run("undo more than recorded never goes negative", func(t *testing.T) {
		tr := NewPositionTracker()
		tr.Record(42)
		tr.Undo(42)
		tr.Undo(42) // now over-undoing
		tr.Undo(42) // and again
		if got := tr.Count(42); got != 0 {
			t.Errorf("Count(42) after over-undo = %d, want 0", got)
		}
	})

	// =========================================================================
	// Round-trip — Record then Undo returns to the original state.
	// =========================================================================

	t.Run("record then undo returns the count to its previous value", func(t *testing.T) {
		tr := NewPositionTracker()
		tr.Record(42)
		tr.Record(42)
		original := tr.Count(42) // 2

		tr.Record(42)
		if got := tr.Count(42); got != original+1 {
			t.Errorf("Count after record = %d, want %d", got, original+1)
		}
		tr.Undo(42)
		if got := tr.Count(42); got != original {
			t.Errorf("Count after undo = %d, want %d", got, original)
		}
	})

	// =========================================================================
	// Multi-position scenario — simulates a real game with repetitions.
	// =========================================================================

	t.Run("a sequence of records and undos tracks counts correctly", func(t *testing.T) {
		tr := NewPositionTracker()

		// Position A appears on move 1, 5, and 9 (threefold repetition).
		// Position B appears on move 2 and 6.
		// Position C appears on move 3 and 7.
		hashA := uint64(100)
		hashB := uint64(200)
		hashC := uint64(300)

		// Forward: A B C A B C A
		tr.Record(hashA)
		tr.Record(hashB)
		tr.Record(hashC)
		tr.Record(hashA) // A count = 2
		tr.Record(hashB)
		tr.Record(hashC)
		tr.Record(hashA) // A count = 3 → threefold!

		if got := tr.Count(hashA); got != 3 {
			t.Errorf("Count(A) = %d, want 3 (threefold)", got)
		}
		if got := tr.Count(hashB); got != 2 {
			t.Errorf("Count(B) = %d, want 2", got)
		}
		if got := tr.Count(hashC); got != 2 {
			t.Errorf("Count(C) = %d, want 2", got)
		}

		// Undo the last 3 moves (A, C, B).
		tr.Undo(hashA) // A count = 2
		tr.Undo(hashC) // C count = 1
		tr.Undo(hashB) // B count = 1

		if got := tr.Count(hashA); got != 2 {
			t.Errorf("Count(A) after undo = %d, want 2", got)
		}
		if got := tr.Count(hashB); got != 1 {
			t.Errorf("Count(B) after undo = %d, want 1", got)
		}
		if got := tr.Count(hashC); got != 1 {
			t.Errorf("Count(C) after undo = %d, want 1", got)
		}
	})

	t.Run("a full game cycle: record all, undo all, counts return to 0", func(t *testing.T) {
		tr := NewPositionTracker()
		hashes := []uint64{1, 2, 3, 1, 2, 1}

		for _, h := range hashes {
			tr.Record(h)
		}

		// Undo in reverse order.
		for i := len(hashes) - 1; i >= 0; i-- {
			tr.Undo(hashes[i])
		}

		for _, h := range hashes {
			if got := tr.Count(h); got != 0 {
				t.Errorf("Count(%d) after full undo = %d, want 0", h, got)
			}
		}
	})

	// =========================================================================
	// Edge cases.
	// =========================================================================

	t.Run("hash 0 works like any other hash", func(t *testing.T) {
		tr := NewPositionTracker()
		tr.Record(0)
		tr.Record(0)
		if got := tr.Count(0); got != 2 {
			t.Errorf("Count(0) = %d, want 2", got)
		}
		tr.Undo(0)
		if got := tr.Count(0); got != 1 {
			t.Errorf("Count(0) after undo = %d, want 1", got)
		}
	})

	t.Run("undo deletes the map entry when count reaches 0 (keeps the map small)", func(t *testing.T) {
		tr := NewPositionTracker()
		tr.Record(42)
		tr.Undo(42)
		// The entry should be gone from the internal map.
		_, exists := tr.counter[42]
		if exists {
			t.Errorf("map entry for 42 should have been deleted after count reached 0")
		}
	})
}

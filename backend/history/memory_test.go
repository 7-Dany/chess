package history

import (
	"testing"

	"github.com/7-Dany/chess/core"
)

// TestMemoryStore verifies that the in-memory history store behaves as a
// correct LIFO stack — push, pop, peek, len, and all.
func TestMemoryStore(t *testing.T) {
	// Helper: build a snapshot with a from/to position for easy identification.
	snap := func(from, to core.Position) core.Snapshot {
		return core.Snapshot{
			Move: core.Move{From: from, To: to},
		}
	}

	// =========================================================================
	// Empty store behavior.
	// =========================================================================

	t.Run("a fresh store has length 0", func(t *testing.T) {
		s := NewMemoryStore()
		if got := s.Len(); got != 0 {
			t.Errorf("Len() = %d, want 0", got)
		}
	})

	t.Run("pop on an empty store returns false", func(t *testing.T) {
		s := NewMemoryStore()
		_, ok := s.Pop()
		if ok {
			t.Errorf("Pop() on empty store should return false")
		}
	})

	t.Run("peek on an empty store returns false", func(t *testing.T) {
		s := NewMemoryStore()
		_, ok := s.Peek()
		if ok {
			t.Errorf("Peek() on empty store should return false")
		}
	})

	t.Run("all on an empty store returns an empty slice", func(t *testing.T) {
		s := NewMemoryStore()
		got := s.All()
		if len(got) != 0 {
			t.Errorf("All() = %v, want empty", got)
		}
	})

	// =========================================================================
	// Push / Pop — LIFO (last in, first out).
	// =========================================================================

	t.Run("push then pop returns the same entry (LIFO)", func(t *testing.T) {
		s := NewMemoryStore()
		entry := snap(core.A1, core.A2)
		s.Push(entry)

		got, ok := s.Pop()
		if !ok {
			t.Fatalf("Pop() returned false, want true")
		}
		if got.Move.From != core.A1 || got.Move.To != core.A2 {
			t.Errorf("Pop() = %v→%v, want A1→A2", got.Move.From, got.Move.To)
		}
	})

	t.Run("push three entries then pop returns them in reverse order", func(t *testing.T) {
		s := NewMemoryStore()
		s.Push(snap(core.A1, core.A2)) // bottom
		s.Push(snap(core.B1, core.B3)) // middle
		s.Push(snap(core.C1, core.C2)) // top

		got1, ok1 := s.Pop()
		got2, ok2 := s.Pop()
		got3, ok3 := s.Pop()

		if !ok1 || !ok2 || !ok3 {
			t.Fatalf("expected three successful pops")
		}
		if got1.Move.To != core.C2 {
			t.Errorf("first pop = %v, want C2 (last pushed)", got1.Move.To)
		}
		if got2.Move.To != core.B3 {
			t.Errorf("second pop = %v, want B3", got2.Move.To)
		}
		if got3.Move.To != core.A2 {
			t.Errorf("third pop = %v, want A2 (first pushed)", got3.Move.To)
		}
	})

	// =========================================================================
	// Len — tracks the count correctly.
	// =========================================================================

	t.Run("Len reflects the number of entries after pushes and pops", func(t *testing.T) {
		s := NewMemoryStore()

		s.Push(snap(core.A1, core.A2))
		if got := s.Len(); got != 1 {
			t.Errorf("Len() after 1 push = %d, want 1", got)
		}

		s.Push(snap(core.B1, core.B3))
		if got := s.Len(); got != 2 {
			t.Errorf("Len() after 2 pushes = %d, want 2", got)
		}

		s.Pop()
		if got := s.Len(); got != 1 {
			t.Errorf("Len() after 1 pop = %d, want 1", got)
		}

		s.Pop()
		if got := s.Len(); got != 0 {
			t.Errorf("Len() after 2 pops = %d, want 0", got)
		}
	})

	// =========================================================================
	// Peek — look at the top without removing it.
	// =========================================================================

	t.Run("peek returns the top entry without removing it", func(t *testing.T) {
		s := NewMemoryStore()
		s.Push(snap(core.A1, core.A2))
		s.Push(snap(core.B1, core.B3)) // this is the top

		got, ok := s.Peek()
		if !ok {
			t.Fatalf("Peek() returned false")
		}
		if got.Move.To != core.B3 {
			t.Errorf("Peek() = %v, want B3", got.Move.To)
		}
		if s.Len() != 2 {
			t.Errorf("Len() after peek = %d, want 2 (peek should not remove)", s.Len())
		}
	})

	// =========================================================================
	// All — returns a copy, oldest first.
	// =========================================================================

	t.Run("all returns all entries in push order (oldest first)", func(t *testing.T) {
		s := NewMemoryStore()
		s.Push(snap(core.A1, core.A2)) // oldest
		s.Push(snap(core.B1, core.B3))
		s.Push(snap(core.C1, core.C2)) // newest

		all := s.All()
		if len(all) != 3 {
			t.Fatalf("All() returned %d entries, want 3", len(all))
		}
		if all[0].Move.To != core.A2 {
			t.Errorf("All()[0] = %v, want A2 (oldest)", all[0].Move.To)
		}
		if all[1].Move.To != core.B3 {
			t.Errorf("All()[1] = %v, want B3", all[1].Move.To)
		}
		if all[2].Move.To != core.C2 {
			t.Errorf("All()[2] = %v, want C2 (newest)", all[2].Move.To)
		}
	})

	t.Run("all returns a copy — modifying it does not affect the store", func(t *testing.T) {
		s := NewMemoryStore()
		s.Push(snap(core.A1, core.A2))

		all := s.All()
		all[0].Move.To = core.H8 // mutate the copy

		// The store should be unaffected.
		got, _ := s.Peek()
		if got.Move.To != core.A2 {
			t.Errorf("store was modified by All() copy: Peek() = %v, want A2", got.Move.To)
		}
	})

	// =========================================================================
	// Round-trip — push N, pop N, back to empty.
	// =========================================================================

	t.Run("push 5 then pop 5 returns to an empty store", func(t *testing.T) {
		s := NewMemoryStore()
		for i := range 5 {
			s.Push(snap(core.Position(i), core.Position(i+1)))
		}
		if s.Len() != 5 {
			t.Fatalf("Len() = %d, want 5", s.Len())
		}

		for range 5 {
			_, ok := s.Pop()
			if !ok {
				t.Fatalf("Pop() returned false before emptying 5 entries")
			}
		}

		if s.Len() != 0 {
			t.Errorf("Len() after popping all = %d, want 0", s.Len())
		}
		_, ok := s.Pop()
		if ok {
			t.Errorf("Pop() on emptied store should return false")
		}
	})
}

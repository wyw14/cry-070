package masking

import "sync"

type Ledger struct {
	mu      sync.Mutex
	batches map[string]Batch
}

func NewLedger() *Ledger { return &Ledger{batches: map[string]Batch{}} }
func (l *Ledger) Put(key string, b Batch) (Batch, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if old, ok := l.batches[key]; ok {
		return old, false
	}
	l.batches[key] = b.Clone()
	return b, true
}
func (l *Ledger) Get(key string) (Batch, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	b, ok := l.batches[key]
	return b.Clone(), ok
}

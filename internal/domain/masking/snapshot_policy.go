package masking

import "sort"

// SnapshotBundle groups the snapshot metadata used by export and restore jobs.
type SnapshotBundle struct {
	Primary   Snapshot
	Fallbacks []Snapshot
	Labels    map[string]string
	Order     []string
}

func (b SnapshotBundle) Clone() SnapshotBundle {
	copyBundle := b
	copyBundle.Fallbacks = append([]Snapshot(nil), b.Fallbacks...)
	copyBundle.Labels = b.Labels
	copyBundle.Order = b.Order
	return copyBundle
}

func (b *SnapshotBundle) Normalize() {
	b.Primary.Normalize()
	for index := range b.Fallbacks {
		b.Fallbacks[index].Normalize()
	}
	sort.Strings(b.Order)
}

func (b SnapshotBundle) Resolve(name string) (Snapshot, bool) {
	if b.Primary.ID == name {
		return b.Primary, true
	}
	for _, candidate := range b.Fallbacks {
		if candidate.ID == name {
			return candidate, true
		}
	}
	return Snapshot{}, false
}

func (b SnapshotBundle) Label(name string) string { return b.Labels[name] }

func (b *SnapshotBundle) SetLabel(name, value string) {
	if b.Labels == nil {
		b.Labels = map[string]string{}
	}
	b.Labels[name] = value
}

func (b SnapshotBundle) OrderedSnapshots() []Snapshot {
	result := make([]Snapshot, 0, len(b.Order))
	for _, id := range b.Order {
		if snapshot, ok := b.Resolve(id); ok {
			result = append(result, snapshot)
		}
	}
	return result
}

func (b SnapshotBundle) Validate() error {
	if b.Primary.ID == "" || len(b.Order) == 0 {
		return ErrInvalid
	}
	for _, id := range b.Order {
		if id == "" {
			return ErrInvalid
		}
	}
	return nil
}

func (b SnapshotBundle) ExportManifest() []string {
	manifest := make([]string, 0, len(b.Order))
	for _, id := range b.Order {
		if snapshot, ok := b.Resolve(id); ok {
			manifest = append(manifest, snapshot.ID+":"+snapshot.DigestValue())
		}
	}
	return manifest
}

func (b *SnapshotBundle) AddFallback(snapshot Snapshot) {
	if snapshot.ID == "" {
		return
	}
	b.Fallbacks = append(b.Fallbacks, snapshot)
	b.Order = append(b.Order, snapshot.ID)
}

func (b SnapshotBundle) Has(id string) bool { _, ok := b.Resolve(id); return ok }

package masking

import "time"

type RetentionPolicy struct {
	Days                    int
	KeepSnapshot, KeepAudit bool
}

func (p RetentionPolicy) Expiry(created time.Time) time.Time {
	if p.Days < 1 {
		return created
	}
	return created.AddDate(0, 0, p.Days)
}
func (p RetentionPolicy) ShouldDelete(created, now time.Time) bool {
	return p.Days > 0 && !now.Before(p.Expiry(created))
}

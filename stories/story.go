package stories

import "time"

// Story provides the structure for HN story fields
type Story struct {
	Id    int
	Time  time.Time
	Title string
	URL   string
}

package models

type Report struct {
	Batches []Batch
	EventsPerSecondAvg int
	RowCount int
}

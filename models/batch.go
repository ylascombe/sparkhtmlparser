package models

type Batch struct {
	BatchTime string
	InputSize int
	SchedulingDelay int
	ProcessingTime float64
	TotalDelay float64
}

package models

type Batch struct {
	BatchTime string
	InputSize int
	SchedulingDelay int
	ProcessingTime float32
	TotalDelay float32
}

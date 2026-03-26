package pinger

type Result struct {
	Region  Region
	Latency int64
	Status  string
}

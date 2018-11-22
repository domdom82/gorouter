package route

import "time"

type endpointElem struct {
	endpoint           *Endpoint
	index              int
	updated            time.Time
	failedAt           *time.Time
	maxConnsPerBackend int64
}

func (e *endpointElem) isOverloaded() bool {
	if e.maxConnsPerBackend == 0 {
		return false
	}

	return e.endpoint.Stats.NumberConnections.value >= e.maxConnsPerBackend
}

func (e *endpointElem) failed() {
	t := time.Now()
	e.failedAt = &t
}

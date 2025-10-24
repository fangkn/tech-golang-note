package main

import (
	"fmt"
	"time"
)

type TimePair struct {
	start int64
	end   int64
}

func splitTimeRange(start, end int64) []TimePair {

	segments := make([]TimePair, 0)
	tStart := time.Unix(start, 0)
	tEnd := time.Unix(end, 0)
	for t := tStart; t.Before(tEnd); {
		segmentStart := t
		segmentEnd := time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())
		sEnd := segmentEnd.Add(-1 * time.Second)
		if segmentEnd.After(tEnd) {
			segmentEnd = tEnd
			sEnd = tEnd
		}

		segments = append(segments, TimePair{start: segmentStart.Unix(), end: sEnd.Unix()})
		t = segmentEnd
	}

	return segments
}

func main() {
	start := time.Date(2024, 9, 11, 12, 11, 1, 0, time.UTC).Unix()
	end := time.Date(2024, 9, 12, 23, 11, 2, 0, time.UTC).Unix()

	segments := splitTimeRange(start, end)

	fmt.Printf("%v - %v\n\n", time.Unix(start, 0), time.Unix(end, 0))
	for _, segment := range segments {
		fmt.Printf("%v - %v\n", time.Unix(segment.start, 0), time.Unix(segment.end, 0))
	}
}

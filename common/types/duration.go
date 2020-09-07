package types

import (
	"encoding/json"
	"fmt"
)

// Duration represents a duration between start and end
type Duration struct {
	Start Timestamp `json:"s"`
	End   Timestamp `json:"e"`
}

func NewDuration(start Timestamp, end Timestamp) *Duration {
	return &Duration{Start: start, End: end}
}

func (d *Duration) ContainTime(time Timestamp) bool {
	if d.Start > time {
		return false
	}
	return d.End >= time
}

func (d *Duration) Contain(other *Duration) bool {
	if d.Start > other.Start {
		return false
	}
	return d.End >= other.End
}

func (d *Duration) Contained(other *Duration) bool {
	if d.Start < other.Start {
		return false
	}
	return d.End <= other.End
}

func (d *Duration) ToJson() []byte {
	durationJson, err := json.Marshal(d)
	if err != nil {
		fmt.Println(err)
	}
	return durationJson
}

func (d *Duration) ToBytes() []byte {
	durBytes := append(d.Start.Bytes(), d.End.Bytes()...)
	return durBytes
}

func (d *Duration) DeepCopy() *Duration {
	dCopy := &Duration{
		Start: d.Start.DeepCopy(),
		End:   d.End.DeepCopy(),
	}
	return dCopy
}

func (d *Duration) ToString() string {
	return string(d.ToJson())
}

func (d *Duration) ToMap() map[string]interface{} {
	dMap := make(map[string]interface{})
	dMap["s"] = d.Start
	dMap["e"] = d.End
	return dMap
}

func GenerateTestDuration() *Duration {
	start := RandomTimestamp()
	end := RandomTimestamp()
	if start > end {
		start, end = end, start
	}
	d := NewDuration(start, end)
	return d
}

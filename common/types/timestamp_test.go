package types

import (
	"fmt"
	"testing"
)

func TestCurrentTime(t *testing.T) {
	fmt.Println(CurrentTime())
}

func TestTimestamp_ToJson(t *testing.T) {
	ts := CurrentTime()
	tsJson := ts.ToJson()
	fmt.Println(tsJson)
}

func TestJsonToTimestamp(t *testing.T) {
	ts := CurrentTime()
	tsJson := ts.ToJson()
	ts2 := JsonToTimestamp(tsJson)
	if ts != ts2 {
		fmt.Println("convert wrong")
	} else {
		fmt.Println("convert pass")
	}
}

func TestTimestamp_ToUint64(t *testing.T) {
	ts := CurrentTime()
	tsInt := ts.Uint64()
	fmt.Println(tsInt)
}

func TestTimestamp_DeepCopy(t *testing.T) {
	ts := CurrentTime()
	ts2 := ts.DeepCopy()
	fmt.Println(ts, ts2)
	ts += 100
	fmt.Println(ts, ts2)
}

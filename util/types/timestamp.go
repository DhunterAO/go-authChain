package types

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Timestamp uint64

const MaxTimestamp = Timestamp(MaxUint64)

func CurrentTime() Timestamp {
	return Timestamp(time.Now().Unix())
}

func (ts Timestamp) Uint64() uint64 {
	return uint64(ts)
}

func TimestampFromUint64(t uint64) Timestamp {
	timestamp := Timestamp(t)
	return timestamp
}

func (ts Timestamp) Bytes() []byte {
	return []byte(strconv.FormatUint(ts.Uint64(), 10))
}

func BytesToTimestamp(t []byte) Timestamp {
	intT, err := strconv.Atoi(string(t))
	if err != nil {
		fmt.Println(err)
	}
	timestamp := Timestamp(intT)
	return timestamp
}

func (ts *Timestamp) ToJson() []byte {
	tsJson, err := json.Marshal(ts)
	if err != nil {
		fmt.Println(err)
	}
	return tsJson
}

func JsonToTimestamp(j []byte) Timestamp {
	var ts Timestamp
	err := json.Unmarshal(j, &ts)
	if err != nil {
		fmt.Println(err)
	}
	return ts
}

func (ts Timestamp) DeepCopy() Timestamp {
	return ts
}

func RandomTimestamp() Timestamp {
	return TimestampFromUint64(rand.Uint64())
}

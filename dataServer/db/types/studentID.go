package types

import (
	util2 "goauth/util"
	"math/rand"
)

const StudentIDLength = 10

type StudentID [StudentIDLength]byte

func (stuID *StudentID) ToHex() string {
	return string(util2.Expand(stuID[:]))
}

func (stuID *StudentID) ToBytes() []byte {
	return stuID[:]
}

func (stuID *StudentID) DeepCopy() StudentID {
	var newStuID StudentID
	copy(newStuID[:], stuID[:])
	return newStuID
}

func RandStudentID() StudentID {
	var id StudentID
	for i := 0; i < StudentIDLength; i += 1 {
		id[i] = byte(rand.Int() % 10)
	}
	return id
}

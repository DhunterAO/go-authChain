package types

import (
	"github.com/goAuthChain/crypto/ecdsa"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"goauth/crypto/sha3"
	"goauth/rlp"
	"goauth/util/hexutil"
	"goauth/util/types"
	"math/big"
	"reflect"
)

/////////// Address
const AddressLength = 20

// Address represents the 20 byte address of an Ethereum account.
type Address [AddressLength]byte

var addressT = reflect.TypeOf(Address{})

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) *Address {
	var a Address
	a.SetBytes(b)
	return &a
}

// BigToAddress returns Address with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func BigToAddress(b *big.Int) *Address { return BytesToAddress(b.Bytes()) }

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) *Address { return BytesToAddress(types.FromHex(s)) }

// IsHexAddress verifies whether a string can represent a valid hex-encoded
// Ethereum address or not.
func IsHexAddress(s string) bool {
	if types.HasHexPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*AddressLength && types.IsHex(s)
}

// Bytes gets the string representation of the underlying address.
func (a Address) Bytes() []byte { return a[:] }

// Big converts an address to a big integer.
func (a Address) Big() *big.Int { return new(big.Int).SetBytes(a[:]) }

// Hash converts an address to a hash by left-padding it with zeros.
func (a Address) Hash() Hash { return BytesToHash(a[:]) }

// Hex returns an EIP55-compliant hex string representation of the address.
func (a Address) Hex() string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return "0x" + string(result)
}

// String implements fmt.Stringer.
func (a Address) String() string {
	return a.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (a Address) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), a[:])
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a) it will panic.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// MarshalText returns the hex representation of a.
func (a Address) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *Address) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Address", input, a[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *Address) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(addressT, input, a[:])
}

// Scan implements Scanner for database/sql.
func (a *Address) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Address", src)
	}
	if len(srcB) != AddressLength {
		return fmt.Errorf("can't scan []byte of len %d into Address, want %d", len(srcB), AddressLength)
	}
	copy(a[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (a Address) Value() (driver.Value, error) {
	return a[:], nil
}

// RandHash generates random hash value used for test
func RandAddress() *Address {
	return BytesToAddress(RandHash().Bytes())
}

// CreateAddress creates an ethereum address given the bytes and the nonce
func CreateAddress(b Address, nonce uint64) *Address {
	data, _ := rlp.EncodeToBytes([]interface{}{b, nonce})
	return BytesToAddress(Keccak256(data)[12:])
}

// CreateAddress2 creates an ethereum address given the address bytes, initial
// contract code hash and a salt.
func CreateAddress2(b Address, salt [32]byte, inithash []byte) *Address {
	return BytesToAddress(Keccak256([]byte{0xff}, b.Bytes(), salt[:], inithash)[12:])
}

func EcdsaPubkeyToAddress(p ecdsa.PublicKey) *Address {
	pubBytes := FromECDSAPub(&p)
	return BytesToAddress(Keccak256(pubBytes[1:])[12:])
}
package utils

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

// Float64s proto: double
type Float64s []float64

// MarshalLogArray conforms to zap ArrayMarshaler
func (nums Float64s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendFloat64(nums[i])
	}
	return nil
}

// Float32s proto: float
type Float32s []float32

// MarshalLogArray conforms to zap ArrayMarshaler
func (nums Float32s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendFloat32(nums[i])
	}
	return nil
}

// Int32s proto: int32 / sint32 / sfixed32
type Int32s []int32

// MarshalLogArray conforms to zap ArrayMarshaler
func (nums Int32s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendInt32(nums[i])
	}
	return nil
}

// Int64s proto: int64 / sint64 / sfixed64
type Int64s []int64

// MarshalLogArray conforms to zap ArrayMarshaler
func (nums Int64s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendInt64(nums[i])
	}
	return nil
}

// Uint32s proto: uint32 / fixed32
type Uint32s []uint32

// MarshalLogArray conforms to zap ArrayMarshaler
func (nums Uint32s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendUint32(nums[i])
	}
	return nil
}

// Uint64s proto: uint64 / fixed64
type Uint64s []uint64

// MarshalLogArray conforms to zap ArrayMarshaler
func (nums Uint64s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendUint64(nums[i])
	}
	return nil
}

// Bools proto: bool
type Bools []bool

// MarshalLogArray conforms to zap ArrayMarshaler
func (bs Bools) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range bs {
		arr.AppendBool(bs[i])
	}
	return nil
}

// StringArray proto: string
type StringArray []string

// MarshalLogArray conforms to zap ArrayMarshaler
func (ss StringArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range ss {
		arr.AppendString(ss[i])
	}
	return nil
}

// ByteStringsArray proto: bytes
type ByteStringsArray [][]byte

// MarshalLogArray conforms to zap ArrayMarshaler
func (bss ByteStringsArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range bss {
		arr.AppendByteString(bss[i])
	}
	return nil
}

// Objects proto: message
type Objects []zapcore.ObjectMarshaler

// MarshalLogArray conforms to zap ArrayMarshaler
func (os Objects) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range os {
		arr.AppendObject(os[i])
	}
	return nil
}

// Interfaces for slower but working types
type Interfaces []interface{}

// MarshalLogArray conforms to zap ArrayMarshaler
func (is Interfaces) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range is {
		arr.AppendReflected(is[i])
	}
	return nil
}

// Stringers proto: enum
type Stringers []fmt.Stringer

// MarshalLogArray conforms to zap ArrayMarshaler
func (ss Stringers) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range ss {
		arr.AppendString(ss[i].String())
	}
	return nil
}

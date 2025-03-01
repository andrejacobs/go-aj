package ajio_test

import (
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
	"testing"

	"github.com/andrejacobs/go-aj/ajio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	u8 := ajio.NewVariableDataUint8()
	assert.Equal(t, 1, u8.PrefixSize())
	assert.Equal(t, math.MaxUint8, int(u8.MaxSize()))
	assert.Equal(t, binary.LittleEndian, u8.ByteOrder())

	u16 := ajio.NewVariableDataUint16()
	assert.Equal(t, 2, u16.PrefixSize())
	assert.Equal(t, math.MaxUint16, int(u16.MaxSize()))
	assert.Equal(t, binary.LittleEndian, u16.ByteOrder())

	u32 := ajio.NewVariableDataUint32()
	assert.Equal(t, 4, u32.PrefixSize())
	assert.Equal(t, math.MaxUint32, int(u32.MaxSize()))
	assert.Equal(t, binary.LittleEndian, u32.ByteOrder())

	u64 := ajio.NewVariableDataUint64()
	assert.Equal(t, 8, u64.PrefixSize())
	assert.Equal(t, uint64(math.MaxUint64), u64.MaxSize())
	assert.Equal(t, binary.LittleEndian, u64.ByteOrder())

	u16BigEndian := ajio.NewVariableDataUint16().BigEndian()
	assert.Equal(t, 2, u16BigEndian.PrefixSize())
	assert.Equal(t, math.MaxUint16, int(u16BigEndian.MaxSize()))
	assert.Equal(t, binary.BigEndian, u16BigEndian.ByteOrder())

	assert.Equal(t, binary.LittleEndian, ajio.NewVariableDataUint8().LittleEndian().ByteOrder())
	assert.Equal(t, binary.BigEndian, ajio.NewVariableDataUint8().BigEndian().ByteOrder())
	assert.Equal(t, binary.NativeEndian, ajio.NewVariableDataUint8().NativeEndian().ByteOrder())

	vd := ajio.NewVariableData()
	assert.Equal(t, binary.LittleEndian, vd.ByteOrder())
}

func TestWriteAndReadUint8(t *testing.T) {
	expectedData := []byte("The quick brown fox")
	buffer := bytes.Buffer{}

	v := ajio.NewVariableDataUint8()
	wcount, err := v.Write(&buffer, expectedData)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+v.PrefixSize(), wcount)

	data, rcount, err := v.Read(&buffer, nil)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+v.PrefixSize(), rcount)
	assert.Equal(t, expectedData, data)
}

func TestWriteAndReadUint16(t *testing.T) {
	expectedData := []byte("The quick brown fox")
	buffer := bytes.Buffer{}

	v := ajio.NewVariableDataUint16()
	wcount, err := v.Write(&buffer, expectedData)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+v.PrefixSize(), wcount)

	data, rcount, err := v.Read(&buffer, nil)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+v.PrefixSize(), rcount)
	assert.Equal(t, expectedData, data)
}

func TestWriteAndReadUint32(t *testing.T) {
	expectedData := []byte("The quick brown fox")
	buffer := bytes.Buffer{}

	v := ajio.NewVariableDataUint32()
	wcount, err := v.Write(&buffer, expectedData)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+v.PrefixSize(), wcount)

	data, rcount, err := v.Read(&buffer, nil)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+v.PrefixSize(), rcount)
	assert.Equal(t, expectedData, data)
}

func TestWriteAndReadUint64(t *testing.T) {
	expectedData := []byte("The quick brown fox")
	buffer := bytes.Buffer{}

	v := ajio.NewVariableDataUint64()
	wcount, err := v.Write(&buffer, expectedData)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+v.PrefixSize(), wcount)

	data, rcount, err := v.Read(&buffer, nil)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+v.PrefixSize(), rcount)
	assert.Equal(t, expectedData, data)
}

func TestReadUsingExistingBuffer(t *testing.T) {
	expectedData := []byte("The quick brown fox")
	buffer := bytes.Buffer{}

	v := ajio.NewVariableDataUint8()
	wcount, err := v.Write(&buffer, expectedData)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+v.PrefixSize(), wcount)

	// Large enough buffer, so no alloc
	intoBuffer := make([]byte, len(expectedData)+4)
	data, rcount, err := v.Read(&buffer, intoBuffer)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+v.PrefixSize(), rcount)
	assert.Equal(t, expectedData, data)
	assert.True(t, samePointer(intoBuffer, data))

	// Smaller buffer, so need an alloc
	buffer.Reset()
	_, err = v.Write(&buffer, expectedData)
	require.NoError(t, err)

	smallBuffer := make([]byte, len(expectedData)-2)
	data, rcount, err = v.Read(&buffer, smallBuffer)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+v.PrefixSize(), rcount)
	assert.Equal(t, expectedData, data)
	assert.False(t, samePointer(smallBuffer, data))

}

func TestWriteTooBig(t *testing.T) {
	tooBig := make([]byte, math.MaxUint8+1)
	buffer := bytes.Buffer{}

	v := ajio.NewVariableDataUint8()
	wcount, err := v.Write(&buffer, tooBig)
	require.Error(t, err)
	assert.Equal(t, 0, wcount)
}

func TestWriteAndReadString(t *testing.T) {
	expected := "The quick brown fox jumped over the lazy dog!"
	buffer := bytes.Buffer{}

	v8 := ajio.NewVariableDataUint8()
	wcount, err := v8.WriteString(&buffer, expected)
	require.NoError(t, err)
	assert.Equal(t, len(expected)+v8.PrefixSize(), wcount)

	s, rcount, err := v8.ReadString(&buffer)
	require.NoError(t, err)
	assert.Equal(t, expected, s)
	assert.Equal(t, len(expected)+v8.PrefixSize(), rcount)

	v16 := ajio.NewVariableDataUint16()
	wcount, err = v16.WriteString(&buffer, expected)
	require.NoError(t, err)
	assert.Equal(t, len(expected)+v16.PrefixSize(), wcount)

	s, rcount, err = v16.ReadString(&buffer)
	require.NoError(t, err)
	assert.Equal(t, expected, s)
	assert.Equal(t, len(expected)+v16.PrefixSize(), rcount)

	v32 := ajio.NewVariableDataUint32()
	wcount, err = v32.WriteString(&buffer, expected)
	require.NoError(t, err)
	assert.Equal(t, len(expected)+v32.PrefixSize(), wcount)

	s, rcount, err = v32.ReadString(&buffer)
	require.NoError(t, err)
	assert.Equal(t, expected, s)
	assert.Equal(t, len(expected)+v32.PrefixSize(), rcount)

	v64 := ajio.NewVariableDataUint64()
	wcount, err = v64.WriteString(&buffer, expected)
	require.NoError(t, err)
	assert.Equal(t, len(expected)+v64.PrefixSize(), wcount)

	s, rcount, err = v64.ReadString(&buffer)
	require.NoError(t, err)
	assert.Equal(t, expected, s)
	assert.Equal(t, len(expected)+v64.PrefixSize(), rcount)
}

//-----------------------------------------------------------------------------
//varint

func TestWriteAndReadVarInt(t *testing.T) {
	expectedData := []byte("The quick brown fox")
	buffer := bytes.Buffer{}

	v := ajio.NewVariableData()
	wcount, err := v.Write(&buffer, expectedData)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+1, wcount)

	data, rcount, err := v.Read(&buffer, nil)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+1, rcount)
	assert.Equal(t, expectedData, data)

	// 2 bytes
	buffer.Reset()
	expectedData = make([]byte, 200)
	wcount, err = v.Write(&buffer, expectedData)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+2, wcount)

	data, rcount, err = v.Read(&buffer, nil)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+2, rcount)
	assert.Equal(t, expectedData, data)
}

func TestReadVarIntUsingExistingBuffer(t *testing.T) {
	expectedData := []byte("The quick brown fox")
	buffer := bytes.Buffer{}

	v := ajio.NewVariableData()
	wcount, err := v.Write(&buffer, expectedData)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+1, wcount)

	// Large enough buffer, so no alloc
	intoBuffer := make([]byte, len(expectedData)+binary.MaxVarintLen64)
	data, rcount, err := v.Read(&buffer, intoBuffer)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+1, rcount)
	assert.Equal(t, expectedData, data)
	assert.True(t, samePointer(intoBuffer, data))

	// Smaller buffer, so need an alloc
	buffer.Reset()
	_, err = v.Write(&buffer, expectedData)
	require.NoError(t, err)

	smallBuffer := make([]byte, len(expectedData)-2)
	data, rcount, err = v.Read(&buffer, smallBuffer)
	require.NoError(t, err)
	assert.Equal(t, len(expectedData)+1, rcount)
	assert.Equal(t, expectedData, data)
	assert.False(t, samePointer(smallBuffer, data))
}

func TestWriteAndReadStringVarInt(t *testing.T) {
	expected := "The quick brown fox jumped over the lazy dog!"
	buffer := bytes.Buffer{}

	v := ajio.NewVariableData()
	wcount, err := v.WriteString(&buffer, expected)
	require.NoError(t, err)
	assert.Equal(t, len(expected)+1, wcount)

	s, rcount, err := v.ReadString(&buffer)
	require.NoError(t, err)
	assert.Equal(t, expected, s)
	assert.Equal(t, len(expected)+1, rcount)

	// 2 bytes
	buffer.Reset()
	expected = string(make([]byte, 200))
	wcount, err = v.WriteString(&buffer, expected)
	require.NoError(t, err)
	assert.Equal(t, len(expected)+2, wcount)

	s, rcount, err = v.ReadString(&buffer)
	require.NoError(t, err)
	assert.Equal(t, expected, s)
	assert.Equal(t, len(expected)+2, rcount)
}

// -----------------------------------------------------------------------------
// https://stackoverflow.com/questions/58636694/how-to-know-if-2-go-maps-reference-the-same-data
func samePointer(x, y interface{}) bool {
	return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
}

// Copyright (C) 2018 LambdaCal Inc.

package lz

import (
	"encoding/binary"

	"github.com/pkg/errors"
)

// file location //nas/software/tmp/dev/rob.odb
// file location //nas/software/tmp/dev/d3plot01

// ---   compress ratio(compared with origin file size):
// use rob.odb
// float 254132624 61222900 76305923  (lz)0.24090925059664908 (gzip)0.30026024128252027 (lz in c)0.23979692981094786
// double 254132624 61493407 76305923 (lz)0.24197368300104594 (gzip)0.30026024128252027

// use d3plot01
// float 14311168 20096271 23602880   (lz)0.17580321635765284 (gzip)0.20647921294969185
// double 114311168 20816458 23602880 (lz)0.18210344941974524 (gzip)0.20647921294969185
// ---

// ---   compress speed   -------
// use d3plot01
// lz float 114311168 2.00834519         (lz float) 54.281318541659665
// lz double  114311168 2.136632822      (lz double)51.02216154198908
// gzip 114311168 1.827797696            (gzip)     59.64315702912452

// use rob.odb
// lz float 254132624 4.253635564        (lz float) 56.97708532458719
// lz double 254132624 4.604696583       (lz double)52.633165313104534
// gzip 254132624 3.341635055            (gzip)     72.52729651225381

const (
	constFloatSz  = 4
	constDoubleSz = 8
	constUint64Sz = 8
)

// pad data in size bytes, size must be constFloatSz or constDoubleSz
func padData(data []byte, size int) []byte {
	pad := []byte{0x44, 0x44, 0x44, 0x44, 0x44, 0x44, 0x44, 0x44}

	if len(data)%size != 0 {
		data = append(data, pad[(len(data)%size):size]...)
	}

	return data
}

// compressInSize Compress data in size bytes
func compressInSize(input []byte, size int) ([]byte, error) {
	realSize := len(input)
	input = padData(input, size)

	output := make([]byte, constUint64Sz)

	var bitsShift uint = 2
	if size == constDoubleSz {
		bitsShift = 3
	}

	segSize := len(input) / size
	segment := make([]byte, segSize)

	for i := 0; i < size; i++ {
		for j := 0; j < segSize; j++ {
			segment[j] = input[(j<<bitsShift)+i]
		}

		deflateSeg, err := GoZip(segment)
		if err != nil {
			return nil, errors.Errorf("lz compress error %v", err)
		}
		sizeBuf := make([]byte, constUint64Sz)
		binary.BigEndian.PutUint64(sizeBuf, uint64(len(deflateSeg)))
		output = append(output, sizeBuf...)
		output = append(output, deflateSeg...)
	}

	binary.BigEndian.PutUint64(output, uint64(realSize))
	return output, nil
}

// decompressInSize decompress data (compressed in size bytes)
func decompressInSize(input []byte, size int) ([]byte, error) {
	realSize := binary.BigEndian.Uint64(input)
	curIdx := constUint64Sz
	output := make([]byte, (int(realSize)+size-1)/size*size)

	var bitsShift uint = 2
	if size == constDoubleSz {
		bitsShift = 3
	}

	for i := 0; i < size; i++ {
		segSize := binary.BigEndian.Uint64(input[curIdx:])
		curIdx += constUint64Sz
		realSeg, err := GoUnZip(input[curIdx : curIdx+int(segSize)])
		curIdx += int(segSize)

		if err != nil {
			return nil, errors.Errorf("lz uncompress error %v", err)
		}

		for j := 0; j < len(realSeg); j++ {
			output[(j<<bitsShift)+i] = realSeg[j]
		}
	}
	return output[:realSize], nil
}

// CompressInFloat Compress data in 4 bytes
func CompressInFloat(input []byte) ([]byte, error) {
	return compressInSize(input, constFloatSz)
}

// DecompressInFloat Uncompress data (compressed in 4 bytes)
func DecompressInFloat(input []byte) ([]byte, error) {
	return decompressInSize(input, constFloatSz)
}

// CompressInDouble compress data in 8 bytes
func CompressInDouble(input []byte) ([]byte, error) {
	return compressInSize(input, constDoubleSz)
}

// DecompressInDouble uncompress data (compressed in 8 bytes)
func DecompressInDouble(input []byte) ([]byte, error) {
	return decompressInSize(input, constDoubleSz)
}

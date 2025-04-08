package lz

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

var enableD3plotTest = false

var randomData = []byte{
	38, 164, 179, 238, 231, 23, 52, 240,
	27, 13, 94, 114, 103, 122, 82, 147,
	190, 74, 72, 69, 209, 113, 155, 12,
	252, 222, 142, 94, 219, 141, 45, 3,
	115, 154, 22, 111, 131, 252, 208, 228,
	129, 192, 245, 193, 13, 116, 120, 70,
	191, 179, 89, 205, 152, 73, 184, 107,
	96, 245, 142, 159, 84, 252, 100, 20,
}

var (
	smallFeatureData = []byte{}
	largeFeatureData = []byte{}
	mixData          = []byte{}
)

var smallData = []float64{
	1.30e10, 1.301e10, 1.302e10, 1.303e10, 1.304e10, 1.305e10, 1.306e10, 1.307e10, 1.308e10, 1.309e10,
	1.31e10, 1.311e10, 1.312e10, 1.313e10, 1.314e10, 1.315e10, 1.316e10, 1.317e10, 1.318e10, 1.319e10,
	1.32e10, 1.321e10, 1.322e10, 1.323e10, 1.324e10, 1.325e10, 1.326e10, 1.327e10, 1.328e10, 1.329e10,
	1.33e10, 1.331e10, 1.332e10, 1.333e10, 1.334e10, 1.335e10, 1.336e10, 1.337e10, 1.338e10, 1.339e10,
	1.34e10, 1.341e10, 1.342e10, 1.343e10, 1.344e10, 1.345e10, 1.346e10, 1.347e10, 1.348e10, 1.349e10,
	1.35e10, 1.351e10, 1.352e10, 1.353e10, 1.354e10, 1.355e10, 1.356e10, 1.357e10, 1.358e10, 1.359e10,
	1.36e10, 1.361e10, 1.362e10, 1.363e10, 1.364e10, 1.365e10, 1.366e10, 1.367e10, 1.368e10, 1.369e10,
	1.37e10, 1.371e10, 1.372e10, 1.373e10, 1.374e10, 1.375e10, 1.376e10, 1.377e10, 1.378e10, 1.379e10,
	1.38e10, 1.381e10, 1.382e10, 1.383e10, 1.384e10, 1.385e10, 1.386e10, 1.387e10, 1.388e10, 1.389e10,
	1.39e10, 1.391e10, 1.392e10, 1.393e10, 1.394e10, 1.395e10, 1.396e10, 1.397e10, 1.398e10, 1.399e10,

	1.30e10, 1.301e10, 1.302e10, 1.303e10, 1.304e10, 1.305e10, 1.306e10, 1.307e10, 1.308e10, 1.309e10,
	1.31e10, 1.311e10, 1.312e10, 1.313e10, 1.314e10, 1.315e10, 1.316e10, 1.317e10, 1.318e10, 1.319e10,
	1.32e10, 1.321e10, 1.322e10, 1.323e10, 1.324e10, 1.325e10, 1.326e10, 1.327e10, 1.328e10, 1.329e10,
}

var largeData = []float64{}

var d3plot01Data [][]byte

type testData struct {
	data []byte
	name string
}

var testDataArr = []testData{}

func init() {
	// make featureData
	{
		var arr = smallData
		w := &bytes.Buffer{}
		for _, v := range arr {
			n := math.Float64bits(v)
			if err := binary.Write(w, binary.BigEndian, &n); err != nil {
				panic(err)
			}
		}

		smallFeatureData = w.Bytes()
	}
	{
		f := 1.00e10
		for i := 0; i < 1024*1024; i++ {
			largeData = append(largeData, f+1.00e5*float64(i))
		}

		var arr = largeData
		w := &bytes.Buffer{}
		for _, v := range arr {
			n := math.Float64bits(v)
			if err := binary.Write(w, binary.BigEndian, &n); err != nil {
				panic(err)
			}
		}

		largeFeatureData = w.Bytes()
	}

	// make mixData
	mixData = append(largeFeatureData, randomData...)

	// d3plot5
	if enableD3plotTest {
		files := []string{
			"data/d3plot01",
			"data/d3plot02",
			"data/d3plot03",
			"data/d3plot04",
			"data/d3plot05",
		}
		{
			for _, file := range files {
				fd, _ := os.Open(file)
				data, _ := io.ReadAll(fd)
				d3plot01Data = append(d3plot01Data, data)
				fd.Close()
			}
		}
	}

	// compare
	{
		testDataArr = []testData{
			{
				data: randomData,
				name: "[random]",
			},
			{
				data: smallFeatureData,
				name: "[small]",
			},
			{
				data: largeFeatureData,
				name: "[large]",
			},
			{
				data: mixData,
				name: "[mix]",
			},
		}

		for i, v := range d3plot01Data {
			testDataArr = append(testDataArr, testData{
				data: v,
				name: fmt.Sprintf("[d3plot0%v]", i+1),
			})
		}
	}
}

func byteCompare(bs1, bs2 []byte) bool {
	if len(bs1) != len(bs2) {
		return false
	}

	for i := 0; i < len(bs1); i++ {
		if bs1[i] != bs2[i] {
			return false
		}
	}
	return true
}

func TestLz(t *testing.T) {
	var alg = "lz:"

	for _, item := range testDataArr {
		t.Log(alg, item.name, "origin length:", len(item.data))
		compressed, err := CompressInDouble(item.data)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(alg, item.name, "compressed length:", len(compressed))
		decompressed, err := DecompressInDouble(compressed)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(float64(len(compressed)) * 100 / float64(len(item.data)))

		b := byteCompare(item.data, decompressed)
		if !b {
			t.Fatal(alg, item.name, "not same")
		}
	}
}

func TestGzip(t *testing.T) {
	var alg = "gzip:"

	for _, item := range testDataArr {
		t.Log(alg, item.name, "origin length:", len(item.data))
		compressed, err := GoZip(item.data)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(alg, item.name, "compressed length:", len(compressed))
		decompressed, err := GoUnZip(compressed)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(float64(len(compressed)) * 100 / float64(len(item.data)))
		b := byteCompare(item.data, decompressed)
		if !b {
			t.Fatal(alg, item.name, "not same")
		}
	}
}

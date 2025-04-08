package floatutil

import "strconv"

func NumberToFloatStr[T float32 | float64](number T, prec int) string {
	if number == 0 {
		return "0"
	}

	return strconv.FormatFloat(float64(number), 'f', prec, 64)
}

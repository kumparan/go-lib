package utils

// Float32PointerToFloat64 :nodoc:
func Float32PointerToFloat64(f *float32) float64 {
	if f != nil {
		return float64(*f)
	}
	return float64(0)
}

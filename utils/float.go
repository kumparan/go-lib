package utils

// Float64PointerToFloat64 :nodoc:
func Float64PointerToFloat64(f *float64) float64 {
	if f != nil {
		return *f
	}
	return float64(0)
}

package convert

// IntToInt32 converts int to int32.
func IntToInt32(i int) int32 {
	return int32(i)
}

// Int32ToInt converts int32 to int.
func Int32ToInt(i int32) int {
	return int(i)
}

// Int64ToInt64 is an identity function for int64 (used by goverter).
func Int64ToInt64(i int64) int64 {
	return i
}

package pkg

type byteType interface {
	~int | ~int32 | ~int64 | ~int16 | ~float64 | ~float32
}

func BytesToMB[T byteType](bytes T) float32 {
	kb := bytes / 1024
	mb := kb / 1024

	return float32(mb)
}

package util

type blackHoleWriter int

var BlackHoleWriter blackHoleWriter

func (b blackHoleWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

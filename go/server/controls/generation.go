package controls

var generationChannel chan uint32

func init() {
	generationChannel = make(chan uint32, 5)

	go func() {
		var generation uint32 = 0

		for {
			generationChannel <- generation
			generation++
		}
	}()
}

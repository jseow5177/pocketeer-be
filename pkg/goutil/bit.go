package goutil

func GetSetBitPos(num uint32) []uint32 {
	pos := []uint32{}

	for i := 0; num > 0; i++ {
		if num&1 == 1 {
			pos = append(pos, uint32(i))
		}
		num >>= 1
	}

	return pos
}

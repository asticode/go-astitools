package astibyte

// ToLength forces the length of a []byte
func ToLength(i []byte, rpl byte, length int) []byte {
	if len(i) == length {
		return i
	} else if len(i) > length {
		return i[:length]
	} else {
		for idx := 0; idx <= length-len(i); idx++ {
			i = append(i, rpl)
		}
		return i
	}
}

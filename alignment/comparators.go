package alignment

func Comparator(match int, miss int) func(byte, byte) int{
	return func(char1 byte, char2 byte) int{
		if(char1 == char2){
			return match
		}
		return miss;
	}
}

func MatrixComparator(biosumm_table map[byte](map[byte]int)) func(byte, byte) int{
	return func(char1 byte, char2 byte) int{
		return biosumm_table[char1][char1];
	}
}
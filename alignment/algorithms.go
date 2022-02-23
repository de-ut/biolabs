package alignment

type elem struct{
	score int
	tag byte
}

const(
	DIAG = 0
	UP = 1
	LEFT = 2
)

func insert(s string, c byte) string{
	return string(c)+s;
}

func reverse(s string) string{
	bytes := []byte(s)
	for i := 0; i < len(bytes)/2; i++ {
		bytes[i], bytes[len(bytes)-1-i] = bytes[len(bytes)-1-i], bytes[i]
	}
	return string(bytes)
}

func max2(a int, b int) int{
	if a > b {
		return a
	}
	return b;
}

func max3(a int, b int, c int) int{
	return max2(max2(a,b), c)
}

func min(a int, b int) int{
	if a < b{
		return a
	}
	return b
}

func max_index(arr []int) (index int){
	current_max := arr[0];
	index = 0;
	for i := 1; i < len(arr); i++{
		if(arr[i] > current_max){
			current_max = arr[i];
			index = i;
		}
	}
	return;
}

func NeedlemanWunsch(gap int, seq1 string, seq2 string, comparator func(byte, byte) int) (score int, align1 string, align2 string){
	table := make([][]elem, len(seq1)+1);
	for i := range table {
		table[i] = make([]elem, len(seq2)+1)
	}
	tags := make([][]byte, len(seq1)+1);
	for i := range tags {
		tags[i] = make([]byte, len(seq2)+1)
	}
	for i := 0; i < len(seq1)+1; i++ {
		for j := 0; j < len(seq2)+1; j++ {
			if i == 0 {
				table[0][j] = elem{gap*j, LEFT}
			}else if j == 0 {
				table[i][0] = elem{gap*i, UP}
			}else{
				diag := table[i-1][j-1].score + comparator(seq1[i-1], seq2[j-1])
				up := table[i-1][j].score + gap;
				left := table[i][j-1].score + gap;
				table[i][j].score = max3(diag, up, left)
				if table[i][j].score == diag {
					table[i][j].tag = DIAG;
				} else if table[i][j].score == up {
					table[i][j].tag = UP
				} else{
					table[i][j].tag = LEFT
				}
			}
		}
	}
	score = table[len(seq1)][len(seq2)].score
	align1 = ""
	align2 = ""
	for i, j := len(seq1), len(seq2); i > 0 || j > 0; {
		if table[i][j].tag == DIAG {
			align1 = insert(align1, seq1[i-1])
			align2 = insert(align2, seq2[j-1])
			i--
			j--
		} else if table[i][j].tag == UP {
			align1 = insert(align1, seq1[i-1])
			align2 = insert(align2, '-')
			i--
		} else {
			align1 = insert(align1, '-')
			align2 = insert(align2, seq2[j-1])
			j--
		}
	}
	return;
}

func calcNWScores(gap int, seq1 string, seq2 string, comparator func(byte, byte) int) (scores []int){
	first := make([]int, len(seq2)+1);
	for i := range first{
		first[i] = i*gap
	}
	scores = make([]int, len(seq2)+1)
	for i := 1; i < len(seq1)+1; i++ {
		scores[0] = first[0]+gap
		for j := 1; j < len(seq2)+1; j++ {
			diag := first[j-1] + comparator(seq1[i-1], seq2[j-1])
			up := first[j] + gap;
			left := scores[j-1] + gap;
			scores[j] = max3(diag, up, left)
		}
		copy(first, scores)
	}
	return;
}

func Hirschberg(gap int, seq1 string, seq2 string, comparator func(byte, byte) int) (score int, align1 string, align2 string){
	align1 = ""
	align2 = ""
	score = 0;
	if len(seq1) == 0{
		for i := 0; i < len(seq2); i++ {
			align1 = insert(align1, '-')
			align2 = seq2
		}
		return;
	}
	if len(seq2) == 0{
		for i := 0; i < len(seq1); i++ {
			align1 = seq1
			align2 = insert(align2, '-')
		}
		return;
	}
	if len(seq1) == 1 || len(seq2) == 1{
		return NeedlemanWunsch(gap, seq1, seq2, comparator)
	}
	middle := len(seq1)/2;
	scores1 := calcNWScores(gap, seq1[:middle], seq2, comparator)
	scores2 := calcNWScores(gap, reverse(seq1[middle:]), reverse(seq2), comparator)
	scores := make([]int, 0)
	for i := 0; i < len(scores1) && i < len(scores2); i++{
		scores = append(scores, scores1[i]+scores2[len(scores2)-1-i])
	}
	index := max_index(scores)
	first_score, first_align1, first_align2 := Hirschberg(gap, seq1[:middle], seq2[:index], comparator)
	second_score, second_align1, second_align2 := Hirschberg(gap, seq1[middle:], seq2[index:], comparator)
    align1 = first_align1 + second_align1
    align2 = first_align2 + second_align2
    score = first_score + second_score
	return;
}

func NeedlemanWunschAffine(oGap int, eGap int, seq1 string, seq2 string, comparator func(byte, byte) int) (score int, align1 string, align2 string){
	table := make([][]elem, len(seq1)+1);
	x := make([][]int, len(seq1)+1);
	y := make([][]int, len(seq1)+1);
	for i := 0; i <= len(seq1); i++{
		table[i] = make([]elem, len(seq2)+1)
		x[i] = make([]int, len(seq2)+1)
		y[i] = make([]int, len(seq2)+1)
	}
	inf := 2*oGap + (len(seq1)+len(seq2))*eGap + 1

	table[0][0] = elem{0, DIAG}
	x[0][0] = inf
	y[0][0] = inf
	for i := 1; i <= len(seq1); i++{
		table[i][0] = elem{inf, UP}
		x[i][0] = oGap + (i-1)*eGap
		y[i][0] = inf
	}
	for j := 1; j <= len(seq2); j++{
		table[0][j] = elem{inf, LEFT}
		x[0][j] = inf
		y[0][j] = oGap + (j-1)*eGap
	}

	for i := 1; i <= len(seq1); i++{
		for j := 1; j <= len(seq2); j++{
			x1 := x[i][j-1] + eGap
			x2 := x[i-1][j] + eGap + oGap
			y1 := y[i-1][j] + eGap + oGap
			y2 := y[i-1][j] + eGap
			t1 := table[i][j-1].score + eGap + oGap
			t2 := table[i-1][j].score + eGap + oGap
			t3 := table[i-1][j-1].score + comparator(seq1[i-1], seq2[j-1])

			x[i][j] = max3(x1, y1, t1)
			y[i][j] = max3(x2, y2, t2)
			table[i][j].score = max3(x[i][j], y[i][j], t3)

			if(table[i][j].score == t3){
				table[i][j].tag = DIAG
			}else if(table[i][j].score == y[i][j]){
				table[i][j].tag = UP
			}else{
				table[i][j].tag = LEFT
			}
		}
	}

	score = table[len(seq1)][len(seq2)].score
	align1 = ""
	align2 = ""
	for i, j := len(seq1), len(seq2); i > 0 || j > 0; {
		if table[i][j].tag == DIAG {
			align1 = insert(align1, seq1[i-1])
			align2 = insert(align2, seq2[j-1])
			i--
			j--
		} else if table[i][j].tag == UP {
			align1 = insert(align1, seq1[i-1])
			align2 = insert(align2, '-')
			i--
		} else {
			align1 = insert(align1, '-')
			align2 = insert(align2, seq2[j-1])
			j--
		}
	}
	return;
}
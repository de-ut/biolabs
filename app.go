package main

import (
	"./alignment"
	"strings"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const(
	NEEDLEMAN_WUNSCH = 0
	NEEDLEMAN_WUNSCH_AFFINE = 1
	HIRSCHBERG = 2
)

type results struct {
	algorithm int
	score int
	align1 []string
	align2 []string
}

type cmdArgs struct {
	inputs []string
	output string
	dict   string
	gap    int
	egap   int
	match  int
	miss   int
	nw     bool
	anw    bool
	hirsch  bool
}

func isNumeric(c byte) bool{
	_, err := strconv.ParseInt(string(c), 10, 8)
	return err == nil

}

func parseArgs() (args cmdArgs) {
	args = cmdArgs{make([]string,0),"","",-2,-2,1,-1,false,false,false }
	help :=
`
Command Line Arguments:

  -h --help

	Show help documentation about command line arguments

  -i --input	(REQUIRED)

	Input from file or from 2 files. Usage: -i <file_path1> [<file_path2>]

  -o --output	(default=none)

	Output to file. Usage:  -o <file_path>

  -g --gap	(default=-2)

	Gap and open gap penalty. Usage: -g <value>

  -e --egap	(default=-2)

	Extansion gap penalty for Hirschberg algorithm. Usage: -e <value>

  -M --match	(default=1)

	Score for symbols match. Usage: -M <value>

  -m --miss	(default=-1)

	Score for symbols missmatch. Usage: -m <value>

  -d --dict	(default=none)

	Change to comparasion with matrix from file. Usage: -d <matrix_file>

  -nw		(REQUIRED ONE OF)

	Launch Needleman-Wunsch algorithm		

  -anw		(REQUIRED ONE OF)

	Launch Needleman-Wunsch algorithm with affine penalty

  -hirsh	(REQUIRED ONE OF)

	Launch Hirschberg algorithm
`
	args_strings := os.Args[1:];
	for i := 0; i < len(args_strings); i++{
		switch args_strings[i] {
		case "-h","--help":
			print(help)
			os.Exit(0)
		case "-i","--input":
			for old_i := i; i < old_i+2; i++ {
				if (i+1 < len(args_strings)) && (args_strings[i+1][0] != '-') {
					args.inputs = append(args.inputs, args_strings[i+1])
				}else if i == old_i{
					println("Error: no input files after flag -i. Start with -h to see help.")
				}else{
					break;
				}
			}
		case "-o","--output":
			if (i+1 < len(args_strings)) && (args_strings[i+1][0] != '-') {
				args.output = args_strings[i+1]
				i++
			}else {
				fmt.Println("Error: no output file after flag -o. Programm will show result in terminal. Start with -h to see help.")
			}
		case "-g","--gap":
			if i+1 < len(args_strings){
				value, err := strconv.Atoi(args_strings[i+1])
				if err != nil{
					fmt.Println("Error: bad value after flag -g. Use integer values. Programm will work with default value = -2.")
				}else{
					args.gap = value
					i++
				}
			}else {
				fmt.Println("Error: no value after flag -g. Programm will work with default value = -2. Start with -h to see help.")
			}
		case "-e","--egap":
			if i+1 < len(args_strings){
				value, err := strconv.Atoi(args_strings[i+1])
				if err != nil{
					fmt.Println("Error: bad value after flag -e. Use integer values. Programm will work with default value = -2.")
				}else{
					args.egap = value
					i++
				}
			}else {
				fmt.Println("Error: no value after flag -e. Programm will work with default value = -2. Start with -h to see help.")
			}
			break;
		case "-M","--match":
			if i+1 < len(args_strings){
				value, err := strconv.Atoi(args_strings[i+1])
				if err != nil{
					fmt.Println("Error: bad value after flag -M. Use integer values. Programm will work with default value = 1.")
				}else{
					args.match = value
					i++
				}
			}else {
				fmt.Println("Error: no value after flag -M. Programm will work with default value = 1. Start with -h to see help.")
			}
			break;
		case "-m","--miss":
			if i+1 < len(args_strings){
				value, err := strconv.Atoi(args_strings[i+1])
				if err != nil{
					fmt.Println("Error: bad value after flag -m. Use integer values. Programm will work with default value = -1.")
				}else{
					args.miss = value
					i++
				}
			}else {
				fmt.Println("Error: no value after flag -m. Programm will work with default value = -1. Start with -h to see help.")
			}
			break;
		case "-d","--dict":
			if (i+1 < len(args_strings)) && (args_strings[i+1][0] != '-') {
				args.dict = args_strings[i+1]
				i++
			}else {
				fmt.Println("Error: no matrix file after flag -d. Programm will work with simple comparator. Start with -h to see help.")
			}
			break;
		case "-nw":
			args.nw = true;
			break;
		case "-anw":
			args.anw = true;
			break;
		case "-hirsch":
			args.hirsch = true;
			break;
		default:
			fmt.Printf("Unknown flag %s. Start with -h to see help.\n", args_strings[i])
		}
	}
	if len(args.inputs) == 0{
		fmt.Println("No input files. Please start with -h to see help.")
		os.Exit(0)
	}
	if !(args.nw || args.anw || args.hirsch){
		fmt.Println("Warning. No algorithm has been selected. Please start with -h to see help")
		os.Exit(0)
	}
	return args
}

func scanMatrix(path string) map[byte](map[byte]int){
	hash := make(map[byte](map[byte]int))
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Failed to open matrix file")
		os.Exit(0)
	}
	defer file.Close()

	lines := make([]string, 0)
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		lines = append(lines, strings.TrimSpace(fileScanner.Text()))
	}

	re := regexp.MustCompile(`\s+`)

	header := func(arr []string) []byte{
		chars := make([]byte, 0)
		for _, v := range arr{
			chars = append(chars, v[0])
		}
		return chars;
	}(re.Split((lines[0]), -1))

	if len(lines)-1 != len(header){
		fmt.Println("Wrong matrix");
		os.Exit(0);
	}

	for i := 0; i < len(header); i++{
		hash[header[i]] = make(map[byte]int);
		mass := re.Split(lines[i+1], -1)[1:]
		if len(mass) != len(header){
			fmt.Println("Wrong matrix");
			os.Exit(0);
		}
		for j := 0; j < len(mass); j++{
			m, err := strconv.Atoi(mass[j]);
			if err != nil{
				fmt.Println("Wrong matrix");
				os.Exit(0)
			}
			hash[header[i]][header[j]] = m;
		}
	}
	return hash;
}

func scanInput(input []string) (seq1 string, seq2 string){
	if len(input) == 1 {
		file, err := os.Open(input[0])
		if err != nil {
			fmt.Println("Failed to open input file")
			os.Exit(0)
		}
		defer file.Close()

		lines := make([]string, 0)
		fileScanner := bufio.NewScanner(file)
		for fileScanner.Scan() {
			lines = append(lines, strings.TrimSpace(fileScanner.Text()))
		}
		var seq_number = 1
		for i := range lines{
			if (lines[i]) == "" {
				seq_number++
				continue;
			}else{
				if seq_number == 1 {
					seq1 += lines[i]
				}else if seq_number == 2 {
					seq2 += lines[i]
				}else{
					break;
				}
			}
		}

	}else if len(input) == 2 {
		file1, err := os.Open(input[0])
		if err != nil {
			fmt.Println("Failed to open input file")
			os.Exit(0)
		}
		defer file1.Close()
		fileScanner := bufio.NewScanner(file1)
		for fileScanner.Scan() {
			seq1 += strings.TrimSpace(fileScanner.Text())
		}
		file2, err := os.Open(input[1])
		if err != nil {
			fmt.Println("Failed to open input file")
			os.Exit(0)
		}
		defer file2.Close()
		fileScanner = bufio.NewScanner(file2)
		for fileScanner.Scan() {
			seq2 += strings.TrimSpace(fileScanner.Text())
		}
	}
	return
}

func formatSequence(seq string, count int) (result []string){
	i := 0
	for ; i+count < len(seq); i+=count {
		result = append(result, seq[i:i+count])
	}
	result = append(result, seq[i:])
	return
}

func printResults(output string, seq1 []string, seq2 []string, res []results){
	if(output == ""){
		fmt.Println("Input:")
		for i := range seq1{
			if i == 0 {
				fmt.Printf("SEQ1: %s\n", seq1[i])
			}else{
				fmt.Printf("      %s\n", seq1[i])
			}
		}
		for i := range seq2{
			if i == 0 {
				fmt.Printf("SEQ1: %s\n", seq2[i])
			}else{
				fmt.Printf("      %s\n", seq2[i])
			}
		}
		fmt.Println()

		for _, r := range res {
			switch r.algorithm{
			case NEEDLEMAN_WUNSCH:
				fmt.Println("Needleman-Wunsch Algorithm:");
			case NEEDLEMAN_WUNSCH_AFFINE:
				fmt.Println("Affine Needleman-Wunsch Algorithm:");
			case HIRSCHBERG:
				fmt.Println("Hirschberg Algorithm:");
			}
			fmt.Printf("Score: %d\n", r.score);
			for i := 0; i < len(r.align1); i++ {
				fmt.Printf("SEQ1: %s\n", r.align1[i])
				fmt.Printf("SEQ2: %s\n", r.align2[i])
			}
			fmt.Println()
		}
	} else {
		f, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			fmt.Println("Failed to open output file")
			os.Exit(0)
		}
		defer f.Close()

		f.Truncate(0);
		fmt.Fprintln(f,"Input:")
		for i := range seq1{
			if i == 0 {
				fmt.Fprintf(f,"SEQ1: %s\n", seq1[i])
			}else{
				fmt.Fprintf(f,"      %s\n", seq1[i])
			}
		}
		for i := range seq2{
			if i == 0 {
				fmt.Fprintf(f,"SEQ1: %s\n", seq2[i])
			}else{
				fmt.Fprintf(f,"      %s\n", seq2[i])
			}
		}
		fmt.Fprintln(f)

		for _, r := range res{
			switch r.algorithm{
			case NEEDLEMAN_WUNSCH:
				fmt.Fprintln(f,"Needleman-Wunsch Algorithm:");
			case NEEDLEMAN_WUNSCH_AFFINE:
				fmt.Fprintln(f,"Affine Needleman-Wunsch Algorithm:");
			case HIRSCHBERG:
				fmt.Fprintln(f,"Hirschberg Algorithm:");
			}
			fmt.Fprintf(f, "Score: %d\n", r.score);
			for i := 0; i < len(r.align1); i++ {
				fmt.Fprintf(f, "SEQ1: %s\n", r.align1[i])
				fmt.Fprintf(f, "SEQ2: %s\n", r.align2[i])
			}
			fmt.Fprintln(f)
		}
	}
}

func main(){
	args := parseArgs();
	seq1, seq2 := scanInput(args.inputs)
	var comparator func(byte, byte) int
	if(args.dict == ""){
		comparator = alignment.Comparator(args.match, args.miss)
	}else {
		hash := scanMatrix(args.dict)
		comparator = alignment.MatrixComparator(hash)
	}
	res := make([]results, 0)
	if(args.nw){
		score, align1, align2 := alignment.NeedlemanWunsch(args.gap, seq1, seq2, comparator)
		res = append(res, results{
			NEEDLEMAN_WUNSCH,
			score,
			formatSequence(align1, 50),
			formatSequence(align2, 50),
		})
	}
	if(args.anw){
		score, align1, align2 := alignment.NeedlemanWunschAffine(args.gap, args.egap, seq1, seq2, comparator)
		res = append(res, results{
			NEEDLEMAN_WUNSCH_AFFINE,
			score,
			formatSequence(align1, 50),
			formatSequence(align2, 50),
		})
	}
	if(args.hirsch){
		score, align1, align2 := alignment.Hirschberg(args.gap, seq1, seq2, comparator)
		res = append(res, results{
			HIRSCHBERG,
			score,
			formatSequence(align1, 50),
			formatSequence(align2, 50),
		})
	}
	printResults(args.output, formatSequence(seq1,50), formatSequence(seq2,50), res)
}
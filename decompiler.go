package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

var dataArr []int
var usedData bool
var instructions []string
var ifBreak bool

func main() {
	ifBreak = false
	usedData = false
	for i := range dataArr { //initialize dataArr but it has nothing -> useless code
		dataArr[i] = 0
	}
	// takes care of input files
	var InputFileName *string  // has the InputFileName
	var OutputFileName *string // var to contain output file name

	programCount := 96 // program counter that starts at 96

	InputFileName = flag.String("i", "", "Gets the input file name")
	OutputFileName = flag.String("o", "", "Gets the output file name")

	flag.Parse() // parses the declared flags

	OutputFileName2 := fmt.Sprintf("%v_sim.txt", *OutputFileName)
	*OutputFileName = fmt.Sprintf("%v_dis.txt", *OutputFileName)

	if flag.NArg() != 0 {
		os.Exit(200)
	}

	if *InputFileName == "" { //error handling
		fmt.Println("Specify an input File. -- go run thisFile.go -i input.txt --")
		os.Exit(1) // infile errors range 1-10
	}

	inputFile, err := os.Open(*InputFileName)
	if err != nil { // error handling (cant find/ open file)
		fmt.Println("Error opening input file:", err)
		os.Exit(2)
	}
	defer func(inputFile *os.File) {
		_ = inputFile.Close()
	}(inputFile)

	scanner := bufio.NewScanner(inputFile)

	if *OutputFileName == "" { //
		// if no output file is specified
		for scanner.Scan() {
			line := scanner.Text()
			// process each line as needed
			fmt.Println(line)

		}
	} else {
		outputFile, err := os.Create(*OutputFileName) // creates the output file with specified name
		if err != nil {
			fmt.Println("Error creating output file:", err)
			os.Exit(11) // outfile errors range 11-20
		}
		defer outputFile.Close()

		writer := bufio.NewWriter(outputFile)

		for scanner.Scan() {
			line := scanner.Text() // reads each line until eof
			allZeros := true
			for _, char := range line {
				if char != '0' {
					allZeros = false
					break // exit the loop as soon as a non-zero character is found
				}
			}
			if len(line) != 32 {
				fmt.Println("Error: Unknown Instruction") // makes sure lines are 32bit
			} else if allZeros {
				toPrint := fmt.Sprintf("%.32s", line[0:])
				toPrint = fmt.Sprintf("%v\t%v\tNOP", toPrint, programCount)
				_, err := writer.WriteString(toPrint + "\n")
				if err != nil {
					fmt.Println("Error writing to output file:", err)
					os.Exit(12) // outfile errors range 11-20
				}
			} else {
				toPrint := toASM(line, programCount) // processes each binary line
				programCount += 4                    // increments PC by 4 each iteration
				_, err := writer.WriteString(toPrint + "\n")
				if err != nil {
					fmt.Println("Error writing to output file:", err)
					os.Exit(12) // outfile errors range 11-20
				}
			}
		}

		if err := writer.Flush(); err != nil {
			fmt.Println("Error flushing writer:", err)
			os.Exit(13) //outfile errors range 11-20
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input file:", err)
		os.Exit(1)
	}

	inputFile, err = os.Open(*OutputFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()
	scanner = bufio.NewScanner(inputFile)
	// open output file
	outputFile, err := os.Create(OutputFileName2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)

	var stateArr [32]int //initializes all the registers to zero
	for i := 0; i < 32; i++ {
		stateArr[i] = 0
	}

	currentInstruction := 0 //keeps track of the current instruction being processed
	cycle := 1              //keeps track of cycle to be incremented
	breakIdx := 0
	breakCounter := 0
	for scanner.Scan() {
		line := scanner.Text() // reads each line until eof
		if line != "" {
			instructions = append(instructions, line)
			checker := strings.Fields(instructions[len(instructions)-1])
			breakCounter += 1
			if checker[0] == "11111110" {
				breakIdx = breakCounter
			}
		}
	}
	breakCounter = 0

	if breakIdx >= len(instructions) {
		checker := strings.Fields(instructions[breakIdx-1])
		src1, _ := strconv.Atoi(checker[6])
		dataArr = append(dataArr, src1+4)
	} else {
		for i := breakIdx; i < len(instructions); i++ {
			checker := strings.Fields(instructions[i])
			src1, _ := strconv.Atoi(checker[1])
			src2, _ := strconv.Atoi(checker[2])

			if breakCounter == 0 || breakCounter%8 == 0 {
				dataArr = append(dataArr, src1, src2)
			} else {
				dataArr = append(dataArr, src2)
			}
			breakCounter++
		}
		usedData = true
	}

	for len(dataArr)%9 != 0 {
		dataArr = append(dataArr, 0)
	}

	var toPrint string
	for ifBreak == false { //while NOT BREAK
		toPrint = processLine(instructions[currentInstruction], &stateArr, cycle, &currentInstruction)
		currentInstruction++
		cycle++
		_, err := writer.WriteString(toPrint + "\n")
		if err != nil {
			fmt.Println("Error writing to output file:", err)
			os.Exit(12) // outfile errors range 11-20
		}
	}
	err = writer.Flush() //flush to output the output
	if err != nil {
		return
	}
}

/*
	 processLine: this function takes the instruction and current instruction, state of the array, and the cycle.
		Contains both the codeMatrix, which is an array containing the opcode bits for each instruction. The instMatrix
		contains the string instructions within a separate array. In this function, the current instruction from main()
		goes through this function and searches for the pre-existing instruction using the codeMatrix. If a match is
		found, then the index is collected, which is used to find the instruction (in the form of a string).
*/
func processLine(aString string, stateArr *[32]int, cycle int, currentInstruction *int) string {
	var codeMatrix [20]string
	var instMatrix [20]string

	// opcode key
	codeMatrix[0] = "10001010000"
	codeMatrix[1] = "10001011000"
	codeMatrix[2] = "1001000100"
	codeMatrix[3] = "10101010000"
	codeMatrix[4] = "10110100"
	codeMatrix[5] = "10110101"
	codeMatrix[6] = "11001011000"
	codeMatrix[7] = "1101000100"
	codeMatrix[8] = "000101"
	codeMatrix[9] = "11111000010"
	codeMatrix[10] = "11111000000"
	codeMatrix[11] = "110100101"
	codeMatrix[12] = "111100101"
	codeMatrix[13] = "11010011010"
	codeMatrix[14] = "11010011011"
	codeMatrix[15] = "11010011100"
	codeMatrix[16] = "11111110"
	codeMatrix[17] = "11101010000"
	codeMatrix[18] = "111111111111"
	codeMatrix[19] = "00000000000000000000000000000000"

	// instruction key
	instMatrix[0] = "AND"
	instMatrix[1] = "ADD"
	instMatrix[2] = "ADDI"
	instMatrix[3] = "ORR"
	instMatrix[4] = "CBZ"
	instMatrix[5] = "CBNZ"
	instMatrix[6] = "SUB"
	instMatrix[7] = "SUBI"
	instMatrix[8] = "B"
	instMatrix[9] = "LDUR"
	instMatrix[10] = "STUR"
	instMatrix[11] = "MOVZ"
	instMatrix[12] = "MOVK"
	instMatrix[13] = "LSR"
	instMatrix[14] = "LSL"
	instMatrix[15] = "ASR"
	instMatrix[16] = "BREAK"
	instMatrix[17] = "EOR"
	instMatrix[18] = "signedBin"
	instMatrix[19] = "NOP"
	//---------------------------------------------//
	aString = strings.ReplaceAll(aString, ",", "")
	tempArr1 := strings.Fields(aString)
	instruction := ""

	for i := 0; i < len(codeMatrix); i++ {
		if tempArr1[0] == codeMatrix[i] {
			instruction = instMatrix[i]
			break
		}
	}
	progCount := "" //program counter
	rOne := ""      //registers 1-3 for the different inputs in the instructions
	rTwo := ""
	rThree := ""
	toPrint := ""

	// each instruction has a different format (i.e. Immediate format, Branching, R format, etc.) and thus each
	// instruction will have a different format. Each case makes sure to include 'R' and '#' when appropriate
	switch instruction {
	// R format instructions
	case "ADD", "SUB", "AND", "ORR", "EOR", "LSL", "ASR", "LSR":
		progCount = tempArr1[5]
		instruction = tempArr1[6]
		rOne = tempArr1[7]
		rTwo = tempArr1[8]
		rThree = tempArr1[9]
		// formatting of each register
		rOne = strings.ReplaceAll(rOne, "R", "")
		rTwo = strings.ReplaceAll(rTwo, "R", "")
		rThree = strings.ReplaceAll(rThree, "R", "")
		rThree = strings.ReplaceAll(rThree, "#", "")
		// results of the simulate function are stored in toPrint
		toPrint = simulate(progCount, instruction, rOne, rTwo, rThree, stateArr, cycle, currentInstruction)
	// I format instructions
	case "ADDI", "SUBI":
		progCount = tempArr1[4]
		instruction = tempArr1[5]
		rOne = tempArr1[6]
		rTwo = tempArr1[7]
		rThree = tempArr1[8]
		// formatting of each register
		rOne = strings.ReplaceAll(rOne, "R", "")
		rTwo = strings.ReplaceAll(rTwo, "R", "")
		rThree = strings.ReplaceAll(rThree, "#", "")
		// results of the simulate function are stored in toPrint
		toPrint = simulate(progCount, instruction, rOne, rTwo, rThree, stateArr, cycle, currentInstruction)
	// D format instructions
	case "STUR", "LDUR":
		progCount = tempArr1[5]
		instruction = tempArr1[6]
		rOne = tempArr1[7]
		rTwo = tempArr1[8]
		rThree = tempArr1[9]
		// formatting of each register
		rOne = strings.ReplaceAll(rOne, "R", "")
		rTwo = strings.ReplaceAll(rTwo, "R", "")
		rThree = strings.ReplaceAll(rThree, "#", "")
		rTwo = strings.ReplaceAll(rTwo, "[", "")
		rThree = strings.ReplaceAll(rThree, "]", "")
		// results of the simulate function are stored in toPrint
		toPrint = simulate(progCount, instruction, rOne, rTwo, rThree, stateArr, cycle, currentInstruction)
	// Conditional branch instructions
	case "CBZ", "CBNZ":
		progCount = tempArr1[3]
		instruction = tempArr1[4]
		rOne = tempArr1[5]
		rTwo = tempArr1[6]
		// formatting of each register
		rOne = strings.ReplaceAll(rOne, "R", "")
		rTwo = strings.ReplaceAll(rTwo, "#", "")
		// results of the simulate function are stored in toPrint
		toPrint = simulate(progCount, instruction, rOne, rTwo, rThree, stateArr, cycle, currentInstruction)
	// Branch format instructions
	case "B":
		progCount = tempArr1[2]
		instruction = tempArr1[3]
		rOne = tempArr1[4]
		// formatting of each register
		rOne = strings.ReplaceAll(rOne, "#", "")
		// results of the simulate function are stored in toPrint
		toPrint = simulate(progCount, instruction, rOne, rTwo, rThree, stateArr, cycle, currentInstruction)
	// IM format instructions
	case "MOVZ", "MOVK":
		progCount = tempArr1[4]
		instruction = tempArr1[5]
		rOne = tempArr1[6]
		rTwo = tempArr1[7]
		rThree = tempArr1[9]
		// formatting of each register
		rOne = strings.ReplaceAll(rOne, "R", "")
		// results of the simulate function are stored in toPrint
		toPrint = simulate(progCount, instruction, rOne, rTwo, rThree, stateArr, cycle, currentInstruction)
	// BREAK instruction
	case "BREAK":
		progCount = tempArr1[6]
		instruction = tempArr1[7]
		// results of the simulate function are stored in toPrint
		toPrint = simulate(progCount, instruction, rOne, rTwo, rThree, stateArr, cycle, currentInstruction)
	// NOP instruction
	case "NOP":
		progCount = tempArr1[1]
		instruction = tempArr1[2]
		// results of the simulate function are stored in toPrint
		toPrint = simulate(progCount, instruction, rOne, rTwo, rThree, stateArr, cycle, currentInstruction)
	// default case prints "rand" if no appropriate instruction case exists
	default:
		fmt.Println("rand")
	}

	return toPrint
}

/*
simulate: this function takes the program counter, instruction, register(s), state of the array, and the cycle as

	parameters. The function uses the parameters and prints out, in the specified format (i.e. tabs, and borders)
	which will be used to determine the value of 'toPrint' within the 'processLine' function. Each instruction case
	will perform/simulate their appropriate instruction
*/
func simulate(pc string, inst string, r1 string, r2 string, r3 string, stateArr *[32]int, cycle int, currentInstruction *int) string {

	toPrint := ""

	// each case performs a different action (i.e. 'ADD' will perform addition between the values of the registers)
	switch inst {
	// adds the src1 and src2 values of stateArr into the destReg of stateArr
	case "ADD":
		destReg, _ := strconv.Atoi(r1)
		src1, _ := strconv.Atoi(r2)
		src2, _ := strconv.Atoi(r3) //converts everything to an int for easier working

		stateArr[destReg] = stateArr[src1] + stateArr[src2] //does the operation and puts it in register

		//the same for Almost ALL functions >> handles creating the output to be put in the output writer
		toPrint += "====================\n"
		toPrint += fmt.Sprintf("cycle:%v\t%v\t%v\tR%v, R%v, R%v\n", cycle, pc, inst, r1, r2, r3)
		toPrint += fmt.Sprintf("\nregisters:")
		for i := 0; i < 32; i++ {
			if i%8 == 0 {
				toPrint += fmt.Sprintf("\nr%02d:\t", i)
			}
			toPrint += fmt.Sprintf("%v\t", stateArr[i])
		}
		if usedData == false { //checks how many data lines created
			toPrint += fmt.Sprintf("\n\ndata:\n")
		} else if usedData == true { //Created at least ONE data line
			toPrint += fmt.Sprintf("\n\ndata:")
			for i := 0; i < len(dataArr); i++ {
				if i == 0 || i%9 == 0 {
					toPrint += fmt.Sprintf("\n%v:", dataArr[i])
				} else {
					toPrint += fmt.Sprintf("%v\t", dataArr[i])
				}
			}
		}
	// subtracts the src2 stateArr value from the stateArr src1 value and stores it into the destReg of stateArr
	case "SUB":
		destReg, _ := strconv.Atoi(r1)
		src1, _ := strconv.Atoi(r2)
		src2, _ := strconv.Atoi(r3)

		stateArr[destReg] = stateArr[src1] - stateArr[src2]

		toPrint += "====================\n"
		toPrint += fmt.Sprintf("cycle:%v\t%v\t%v\tR%v, R%v, R%v\n", cycle, pc, inst, r1, r2, r3)
		toPrint += fmt.Sprintf("\nregisters:")
		for i := 0; i < 32; i++ {
			if i%8 == 0 {
				toPrint += fmt.Sprintf("\nr%02d:\t", i)
			}
			toPrint += fmt.Sprintf("%v\t", stateArr[i])
		}
		if usedData == false { //checks how many data lines created
			toPrint += fmt.Sprintf("\n\ndata:\n")
		} else if usedData == true { //Created at least ONE data line
			toPrint += fmt.Sprintf("\n\ndata:")
			for i := 0; i < len(dataArr); i++ {
				if i == 0 || i%9 == 0 {
					toPrint += fmt.Sprintf("\n%v:", dataArr[i])
				} else {
					toPrint += fmt.Sprintf("%v\t", dataArr[i])
				}
			}
		}
	// checks to see if the elements of stateArr, src1 and src2, both have a bit value of 1 at the same position and
	// keeps the bit in that location if they do, otherwise the bit value becomes 0
	case "AND":
		destReg, _ := strconv.Atoi(r1)
		src1, _ := strconv.Atoi(r2)
		src2, _ := strconv.Atoi(r3)

		stateArr[destReg] = stateArr[src1] & stateArr[src2]

		toPrint += "====================\n"
		toPrint += fmt.Sprintf("cycle:%v\t%v\t%v\tR%v, R%v, R%v\n", cycle, pc, inst, r1, r2, r3)
		toPrint += fmt.Sprintf("\nregisters:")
		for i := 0; i < 32; i++ {
			if i%8 == 0 {
				toPrint += fmt.Sprintf("\nr%02d:\t", i)
			}
			toPrint += fmt.Sprintf("%v\t", stateArr[i])
		}
		if usedData == false { //checks how many data lines created
			toPrint += fmt.Sprintf("\n\ndata:\n")
		} else if usedData == true { //Created at least ONE data line
			toPrint += fmt.Sprintf("\n\ndata:")
			for i := 0; i < len(dataArr); i++ {
				if i == 0 || i%9 == 0 {
					toPrint += fmt.Sprintf("\n%v:", dataArr[i])
				} else {
					toPrint += fmt.Sprintf("%v\t", dataArr[i])
				}
			}
		}
	// checks to see if the elements of stateArr, src1 and src2, and places a result of 1 bit at the given position
	// if at least one of the elements has a bit value of 1 in that position
	case "ORR":
		destReg, _ := strconv.Atoi(r1)
		src1, _ := strconv.Atoi(r2)
		src2, _ := strconv.Atoi(r3)

		stateArr[destReg] = stateArr[src1] | stateArr[src2]

		toPrint += "====================\n"
		toPrint += fmt.Sprintf("cycle:%v\t%v\t%v\tR%v, R%v, R%v\n", cycle, pc, inst, r1, r2, r3)
		toPrint += fmt.Sprintf("\nregisters:")
		for i := 0; i < 32; i++ {
			if i%8 == 0 {
				toPrint += fmt.Sprintf("\nr%02d:\t", i)
			}
			toPrint += fmt.Sprintf("%v\t", stateArr[i])
		}
		if usedData == false { //checks how many data lines created
			toPrint += fmt.Sprintf("\n\ndata:\n")
		} else if usedData == true { //Created at least ONE data line
			toPrint += fmt.Sprintf("\n\ndata:")
			for i := 0; i < len(dataArr); i++ {
				if i == 0 || i%9 == 0 {
					toPrint += fmt.Sprintf("\n%v:", dataArr[i])
				} else {
					toPrint += fmt.Sprintf("%v\t", dataArr[i])
				}
			}
		}
	// adds the immediate value of src2 and the src1 value of stateArr into the destReg of stateArr
	case "ADDI":
		destReg, _ := strconv.Atoi(r1)
		src1, _ := strconv.Atoi(r2)
		src2, _ := strconv.Atoi(r3)

		stateArr[destReg] = stateArr[src1] + src2

		toPrint += "====================\n"
		toPrint += fmt.Sprintf("cycle:%v\t%v\t%v\tR%v, R%v, #%v\n", cycle, pc, inst, r1, r2, r3)
		toPrint += fmt.Sprintf("\nregisters:")
		for i := 0; i < 32; i++ {
			if i%8 == 0 {
				toPrint += fmt.Sprintf("\nr%02d:\t", i)
			}
			toPrint += fmt.Sprintf("%v\t", stateArr[i])
		}
		if usedData == false { //checks how many data lines created
			toPrint += fmt.Sprintf("\n\ndata:\n")
		} else if usedData == true { //Created at least ONE data line
			toPrint += fmt.Sprintf("\n\ndata:")
			for i := 0; i < len(dataArr); i++ {
				if i == 0 || i%9 == 0 {
					toPrint += fmt.Sprintf("\n%v:", dataArr[i])
				} else {
					toPrint += fmt.Sprintf("%v\t", dataArr[i])
				}
			}
		}

	// subtracts the immediate value of src2 from the src1 value of stateArr into the destReg of stateArr
	case "SUBI":
		destReg, _ := strconv.Atoi(r1)
		src1, _ := strconv.Atoi(r2)
		src2, _ := strconv.Atoi(r3)

		stateArr[destReg] = stateArr[src1] - src2

		toPrint += "====================\n"
		toPrint += fmt.Sprintf("cycle:%v\t%v\t%v\tR%v, R%v, #%v\n", cycle, pc, inst, r1, r2, r3)
		toPrint += fmt.Sprintf("\nregisters:")
		for i := 0; i < 32; i++ {
			if i%8 == 0 {
				toPrint += fmt.Sprintf("\nr%02d:\t", i)
			}
			toPrint += fmt.Sprintf("%v\t", stateArr[i])
		}
		if usedData == false { //checks how many data lines created
			toPrint += fmt.Sprintf("\n\ndata:\n")
		} else if usedData == true { //Created at least ONE data line
			toPrint += fmt.Sprintf("\n\ndata:")
			for i := 0; i < len(dataArr); i++ {
				if i == 0 || i%9 == 0 {
					toPrint += fmt.Sprintf("\n%v:", dataArr[i])
				} else {
					toPrint += fmt.Sprintf("%v\t", dataArr[i])
				}
			}
		}
	// adds the value of stateArr[src1] with (src2 * 4) and places the value into the dataArr at the 0th position, while
	// placing the value of stateArr[destReg] into dataArr at the 1st position
	case "STUR":
		usedData = true
		destReg, _ := strconv.Atoi(r1)
		src1, _ := strconv.Atoi(r2)
		src2, _ := strconv.Atoi(r3)
		idx := 0
		memoLoc := dataArr[0] //212
		iterate := 0

		temp := stateArr[src1] + (src2 * 4) //finds the memory slot for which data will be stored.

		counterVar := 0
		for i := 0; i < len(dataArr); i += 9 {
			iterate += 1
			if temp >= dataArr[i] && temp <= dataArr[i]+28 {
				idx = 1
				for memoLoc != temp {
					counterVar++
					if counterVar%8 == 0 {
						idx += 1
					}
					idx += 1
					memoLoc += 4
				}
				dataArr[idx] = stateArr[destReg]
				break
			} else {
				dataArr = append(dataArr, dataArr[0]+(32*iterate))
				dataArr = append(dataArr, 0, 0, 0, 0, 0, 0, 0, 0)
			}
		}
		toPrint += "====================\n"
		toPrint += fmt.Sprintf("cycle:%v\t%v\t%v\tR%v, [R%v, #%v]\n", cycle, pc, inst, r1, r2, r3)
		toPrint += fmt.Sprintf("\nregisters:")
		for i := 0; i < 32; i++ {
			if i%8 == 0 {
				toPrint += fmt.Sprintf("\nr%02d:\t", i)
			}
			toPrint += fmt.Sprintf("%v\t", stateArr[i])
		}
		if usedData == false { //checks how many data lines created
			toPrint += fmt.Sprintf("\n\ndata:\n")
		} else if usedData == true { //Created at least ONE data line
			toPrint += fmt.Sprintf("\n\ndata:")
			for i := 0; i < len(dataArr); i++ {
				if i == 0 || i%9 == 0 {
					toPrint += fmt.Sprintf("\n%v:", dataArr[i])
				} else {
					toPrint += fmt.Sprintf("%v\t", dataArr[i])
				}
			}
		}
	case "LDUR":
		destReg, _ := strconv.Atoi(r1)
		src1, _ := strconv.Atoi(r2)
		src2, _ := strconv.Atoi(r3)

		idx := 0 //counter
		rand := dataArr[0]
		temp := stateArr[src1] + (src2 * 4) //mem Loc to be looked at
		iterate := 0
		counterVar := 0

		for i := 0; i < len(dataArr); i += 9 {
			iterate += 1
			if temp >= dataArr[i] && temp <= dataArr[i]+28 {
				idx = 1
				if rand != temp {
					counterVar++
					if counterVar%8 == 0 {
						idx += 1
					}
					idx += 1
					rand += 4
				}
				stateArr[destReg] = dataArr[idx]
				break
			} else {
				stateArr[destReg] = 0
			}
		}
		toPrint += "====================\n"
		toPrint += fmt.Sprintf("cycle:%v\t%v\t%v\tR%v, [R%v, #%v]\n", cycle, pc, inst, r1, r2, r3)
		toPrint += fmt.Sprintf("\nregisters:")
		for i := 0; i < 32; i++ {
			if i%8 == 0 {
				toPrint += fmt.Sprintf("\nr%02d:\t", i)
			}
			toPrint += fmt.Sprintf("%v\t", stateArr[i])
		}
		if usedData == false { //checks how many data lines created
			toPrint += fmt.Sprintf("\n\ndata:\n")
		} else if usedData == true { //Created at least ONE data line
			toPrint += fmt.Sprintf("\n\ndata:")
			for i := 0; i < len(dataArr); i++ {
				if i == 0 || i%9 == 0 {
					toPrint += fmt.Sprintf("\n%v:", dataArr[i])
				} else {
					toPrint += fmt.Sprintf("%v\t", dataArr[i])
				}
			}
		}
	// if the registers value is 0 then the current instruction will move forward by the amount in src1
	case "CBZ":
		destReg, _ := strconv.Atoi(r1)
		src1, _ := strconv.Atoi(r2)

		if stateArr[destReg] == 0 {
			*currentInstruction += src1 - 1
		}

		toPrint += "====================\n"
		toPrint += fmt.Sprintf("cycle:%v\t%v\t%v\tR%v, #%v\n", cycle, pc, inst, r1, r2)
		toPrint += fmt.Sprintf("\nregisters:")
		for i := 0; i < 32; i++ {
			if i%8 == 0 {
				toPrint += fmt.Sprintf("\nr%02d:\t", i)
			}
			toPrint += fmt.Sprintf("%v\t", stateArr[i])
		}
		if usedData == false { //checks how many data lines created
			toPrint += fmt.Sprintf("\n\ndata:\n")
		} else if usedData == true { //Created at least ONE data line
			toPrint += fmt.Sprintf("\n\ndata:")
			for i := 0; i < len(dataArr); i++ {
				if i == 0 || i%9 == 0 {
					toPrint += fmt.Sprintf("\n%v:", dataArr[i])
				} else {
					toPrint += fmt.Sprintf("%v\t", dataArr[i])
				}
			}
		}
	// branches the current instruction forward by the amount in the destReg
	case "B":
		destReg, _ := strconv.Atoi(r1)
		*currentInstruction += destReg - 1

		toPrint += "====================\n"
		toPrint += fmt.Sprintf("cycle:%v\t%v\t%v\t#%v\n", cycle, pc, inst, r1)
		toPrint += fmt.Sprintf("\nregisters:")
		for i := 0; i < 32; i++ {
			if i%8 == 0 {
				toPrint += fmt.Sprintf("\nr%02d:\t", i)
			}
			toPrint += fmt.Sprintf("%v\t", stateArr[i])
		}
		if usedData == false { //checks how many data lines created
			toPrint += fmt.Sprintf("\n\ndata:\n")
		} else if usedData == true { //Created at least ONE data line
			toPrint += fmt.Sprintf("\n\ndata:")
			for i := 0; i < len(dataArr); i++ {
				if i == 0 || i%9 == 0 {
					toPrint += fmt.Sprintf("\n%v:", dataArr[i])
				} else {
					toPrint += fmt.Sprintf("%v\t", dataArr[i])
				}
			}
		}
	// if the bits in stateArr[scr1] and stateArr[scr2] are the same then the bit for that position is 0, but if they
	// are different from the bit in that position will be 1
	case "EOR":
		destReg, _ := strconv.Atoi(r1)
		src1, _ := strconv.Atoi(r2)
		src2, _ := strconv.Atoi(r3)

		stateArr[destReg] = stateArr[src1] ^ stateArr[src2]
		toPrint += "====================\n"
		toPrint += fmt.Sprintf("cycle:%v\t%v\t%v\tR%v, R%v, R%v\n", cycle, pc, inst, r1, r2, r3)
		toPrint += fmt.Sprintf("\nregisters:")
		for i := 0; i < 32; i++ {
			if i%8 == 0 {
				toPrint += fmt.Sprintf("\nr%02d:\t", i)
			}
			toPrint += fmt.Sprintf("%v\t", stateArr[i])
		}
		if usedData == false { //checks how many data lines created
			toPrint += fmt.Sprintf("\n\ndata:\n")
		} else if usedData == true { //Created at least ONE data line
			toPrint += fmt.Sprintf("\n\ndata:")
			for i := 0; i < len(dataArr); i++ {
				if i == 0 || i%9 == 0 {
					toPrint += fmt.Sprintf("\n%v:", dataArr[i])
				} else {
					toPrint += fmt.Sprintf("%v\t", dataArr[i])
				}
			}
		}
	// stores the value of stateArr[src1] into stateArr[destReg] when the bits are shifted logically to the left by an
	// amount of src2
	case "LSL":
		destReg, _ := strconv.Atoi(r1)
		src1, _ := strconv.Atoi(r2)
		src2, _ := strconv.Atoi(r3)

		stateArr[destReg] = stateArr[src1] << src2

		toPrint += "====================\n"
		toPrint += fmt.Sprintf("cycle:%v\t%v\t%v\tR%v, R%v, #%v\n", cycle, pc, inst, r1, r2, r3)
		toPrint += fmt.Sprintf("\nregisters:")
		for i := 0; i < 32; i++ {
			if i%8 == 0 {
				toPrint += fmt.Sprintf("\nr%02d:\t", i)
			}
			toPrint += fmt.Sprintf("%v\t", stateArr[i])
		}
		if usedData == false { //checks how many data lines created
			toPrint += fmt.Sprintf("\n\ndata:\n")
		} else if usedData == true { //Created at least ONE data line
			toPrint += fmt.Sprintf("\n\ndata:")
			for i := 0; i < len(dataArr); i++ {
				if i == 0 || i%9 == 0 {
					toPrint += fmt.Sprintf("\n%v:", dataArr[i])
				} else {
					toPrint += fmt.Sprintf("%v\t", dataArr[i])
				}
			}
		}

	// stores the value of stateArr[src1] into stateArr[destReg] when the bits are arithmetically shifted to the right
	// by an amount of src2
	case "ASR":
		destReg, _ := strconv.Atoi(r1)
		src1, _ := strconv.Atoi(r2)
		src2, _ := strconv.Atoi(r3)

		stateArr[destReg] = stateArr[src1] >> src2

		toPrint += "====================\n"
		toPrint += fmt.Sprintf("cycle:%v\t%v\t%v\tR%v, R%v, #%v\n", cycle, pc, inst, r1, r2, r3)
		toPrint += fmt.Sprintf("\nregisters:")
		for i := 0; i < 32; i++ {
			if i%8 == 0 {
				toPrint += fmt.Sprintf("\nr%02d:\t", i)
			}
			toPrint += fmt.Sprintf("%v\t", stateArr[i])
		}
		if usedData == false { //checks how many data lines created
			toPrint += fmt.Sprintf("\n\ndata:\n")
		} else if usedData == true { //Created at least ONE data line
			toPrint += fmt.Sprintf("\n\ndata:")
			for i := 0; i < len(dataArr); i++ {
				if i == 0 || i%9 == 0 {
					toPrint += fmt.Sprintf("\n%v:", dataArr[i])
				} else {
					toPrint += fmt.Sprintf("%v\t", dataArr[i])
				}
			}
		}
	// stores the value of stateArr[src1] into stateArr[destReg] when the bits are shifted logically to the right by an
	// amount of src2
	case "LSR":
		destReg, _ := strconv.Atoi(r1)
		src1, _ := strconv.Atoi(r2)
		src2, _ := strconv.Atoi(r3)

		stateArr[destReg] = int(uint(stateArr[src1]) >> src2 & 0x7FFFFFFF)

		toPrint += "====================\n"
		toPrint += fmt.Sprintf("cycle:%v\t%v\t%v\tR%v, R%v, #%v\n", cycle, pc, inst, r1, r2, r3)
		toPrint += fmt.Sprintf("\nregisters:")
		for i := 0; i < 32; i++ {
			if i%8 == 0 {
				toPrint += fmt.Sprintf("\nr%02d:\t", i)
			}
			toPrint += fmt.Sprintf("%v\t", stateArr[i])
		}
		if usedData == false { //checks how many data lines created
			toPrint += fmt.Sprintf("\n\ndata:\n")
		} else if usedData == true { //Created at least ONE data line
			toPrint += fmt.Sprintf("\n\ndata:")
			for i := 0; i < len(dataArr); i++ {
				if i == 0 || i%9 == 0 {
					toPrint += fmt.Sprintf("\n%v:", dataArr[i])
				} else {
					toPrint += fmt.Sprintf("%v\t", dataArr[i])
				}
			}
		}
	// moves the bits of src1 by the amount in src2, while keeping the original values of src1 in the same position, and
	// stores the new value into the stateArr[destReg]
	case "MOVK":
		destReg, _ := strconv.Atoi(r1)
		src1, _ := strconv.Atoi(r2)
		src2, _ := strconv.Atoi(r3)

		stateArr[destReg] |= src1 << (src2 / 2)

		toPrint += "====================\n"
		toPrint += fmt.Sprintf("cycle:%v\t%v\t%v\tR%v, R%v, R%v\n", cycle, pc, inst, r1, r2, r3)
		toPrint += fmt.Sprintf("\nregisters:")
		for i := 0; i < 32; i++ {
			if i%8 == 0 {
				toPrint += fmt.Sprintf("\nr%02d:\t", i)
			}
			toPrint += fmt.Sprintf("%v\t", stateArr[i])
		}
		if usedData == false { //checks how many data lines created
			toPrint += fmt.Sprintf("\n\ndata:\n")
		} else if usedData == true { //Created at least ONE data line
			toPrint += fmt.Sprintf("\n\ndata:")
			for i := 0; i < len(dataArr); i++ {
				if i == 0 || i%9 == 0 {
					toPrint += fmt.Sprintf("\n%v:", dataArr[i])
				} else {
					toPrint += fmt.Sprintf("%v\t", dataArr[i])
				}
			}
		}
	// moves the bits of src1 by the amount in src2, while making the original values of src1 equal to zero, and stores
	// the new value into the stateArr[destReg]
	case "MOVZ":
		destReg, _ := strconv.Atoi(r1)
		src1, _ := strconv.Atoi(r2)
		src2, _ := strconv.Atoi(r3)

		stateArr[destReg] = 0
		stateArr[destReg] = src1 << (src2 / 2)

		toPrint += "====================\n"
		toPrint += fmt.Sprintf("cycle:%v\t%v\t%v\tR%v, R%v, R%v\n", cycle, pc, inst, r1, r2, r3)
		toPrint += fmt.Sprintf("\nregisters:")
		for i := 0; i < 32; i++ {
			if i%8 == 0 {
				toPrint += fmt.Sprintf("\nr%02d:\t", i)
			}
			toPrint += fmt.Sprintf("%v\t", stateArr[i])
		}
		if usedData == false { //checks how many data lines created
			toPrint += fmt.Sprintf("\n\ndata:\n")
		} else if usedData == true { //Created at least ONE data line
			toPrint += fmt.Sprintf("\n\ndata:")
			for i := 0; i < len(dataArr); i++ {
				if i == 0 || i%9 == 0 {
					toPrint += fmt.Sprintf("\n%v:", dataArr[i])
				} else {
					toPrint += fmt.Sprintf("%v\t", dataArr[i])
				}
			}
		}
	// causes the software to create a breakpoint where the program will stop, and will not generate any results beyond
	// this point
	case "BREAK":
		ifBreak = true
		toPrint += "====================\n"
		toPrint += fmt.Sprintf("cycle:%v\t%v\t%v\n", cycle, pc, inst)
		toPrint += fmt.Sprintf("\nregisters:")
		for i := 0; i < 32; i++ {
			if i%8 == 0 {
				toPrint += fmt.Sprintf("\nr%02d:\t", i)
			}
			toPrint += fmt.Sprintf("%v\t", stateArr[i])
		}
		if usedData == false { //checks how many data lines created
			toPrint += fmt.Sprintf("\n\ndata:\n")
		} else if usedData == true { //Created at least ONE data line
			toPrint += fmt.Sprintf("\n\ndata:")
			for i := 0; i < len(dataArr); i++ {
				if i == 0 || i%9 == 0 {
					toPrint += fmt.Sprintf("\n%v:", dataArr[i])
				} else {
					toPrint += fmt.Sprintf("%v\t", dataArr[i])
				}
			}
		}
	default:
		//
	}

	return toPrint
}

/*
	 toASM: this function takes a string, of bits, and uses certain lengths (6, 8, 9, 10, 11, 12) to identify which
		instruction is being used, which is then compared to pre-established bit strings, along with the values of the
		registers that follow the instruction. Then the function calls upon binaryToDecimal and signedBinToDecimal to
		convert the binary strings to decimal values
*/
func toASM(aString string, progCount int) string {
	var codeMatrix [20]string
	var instMatrix [20]string
	exist := false

	//opcode key
	codeMatrix[0] = "10001010000"
	codeMatrix[1] = "10001011000"
	codeMatrix[2] = "1001000100"
	codeMatrix[3] = "10101010000"
	codeMatrix[4] = "10110100"
	codeMatrix[5] = "10110101"
	codeMatrix[6] = "11001011000"
	codeMatrix[7] = "1101000100"
	codeMatrix[8] = "000101"
	codeMatrix[9] = "11111000010"
	codeMatrix[10] = "11111000000"
	codeMatrix[11] = "110100101"
	codeMatrix[12] = "111100101"
	codeMatrix[13] = "11010011010"
	codeMatrix[14] = "11010011011"
	codeMatrix[15] = "11010011100"
	codeMatrix[16] = "11111110"
	codeMatrix[17] = "11101010000"
	codeMatrix[18] = "111111111111"

	//instruction key
	instMatrix[0] = "AND"
	instMatrix[1] = "ADD"
	instMatrix[2] = "ADDI"
	instMatrix[3] = "ORR"
	instMatrix[4] = "CBZ"
	instMatrix[5] = "CBNZ"
	instMatrix[6] = "SUB"
	instMatrix[7] = "SUBI"
	instMatrix[8] = "B"
	instMatrix[9] = "LDUR"
	instMatrix[10] = "STUR"
	instMatrix[11] = "MOVZ"
	instMatrix[12] = "MOVK"
	instMatrix[13] = "LSR"
	instMatrix[14] = "LSL"
	instMatrix[15] = "ASR"
	instMatrix[16] = "BREAK"
	instMatrix[17] = "EOR"
	instMatrix[18] = "signedBin"

	var temp = ""
	for i := 6; i < 17; i++ {
		temp = aString[0:i]
		for _, j := range codeMatrix {
			if temp == j {
				exist = true
			}
		}
		if exist {
			break
		}
	}
	toPrint := ""
	codeLength := len(temp)
	switch codeLength {
	case 6:
		toPrint = fmt.Sprintf("%.6s %.26s", aString[0:6], aString[6:])
		tempArr1 := strings.Fields(toPrint)
		var tempArr2 [2]string
		for i := 0; i <= 10; i++ {
			if codeMatrix[i] == tempArr1[0] {
				tempArr2[0] = instMatrix[i]
			}
		}
		for i := 1; i < 2; i++ {
			decimal, err := signedBinToDecimal(tempArr1[i])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				tempArr2[i] = decimal.String()
			}
		}
		toPrint = fmt.Sprintf("%v %v\t%v\t#%v", toPrint, progCount, tempArr2[0], tempArr2[1])
		//fmt.Println(toPrint)
	case 8:
		if aString[0:8] == "11111110" {
			toPrint = fmt.Sprintf("%.8s %.3s %.5s %.5s %.5s %.6s", aString[0:8], aString[8:11], aString[11:16], aString[16:21], aString[21:26], aString[26:])
			tempArr1 := strings.Fields(toPrint)
			var tempArr2 [6]string

			tempArr2[0] = "BREAK"
			for i := 1; i < 6; i++ {
				decimal, err := binaryToDecimal(tempArr1[i])
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				} else {
					tempArr2[i] = strconv.FormatInt(decimal, 10)
				}
			}
			toPrint = fmt.Sprintf("%v\t%v\t%v", toPrint, progCount, tempArr2[0])
		} else {
			toPrint = fmt.Sprintf("%.8s %.19s %.5s  ", aString[0:8], aString[8:27], aString[27:])
			tempArr1 := strings.Fields(toPrint)
			var tempArr2 [3]string

			for i := 0; i <= 15; i++ {
				if codeMatrix[i] == tempArr1[0] {
					tempArr2[0] = instMatrix[i]
				}
			}
			for i := 2; i < 3; i++ {
				decimal, err := binaryToDecimal(tempArr1[i])
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				} else {
					tempArr2[i] = strconv.FormatInt(decimal, 10)
				}
			}
			for i := 1; i < 2; i++ {
				decimal, err := signedBinToDecimal(tempArr1[i])
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				} else {
					tempArr2[i] = decimal.String()
				}
			}
			toPrint = fmt.Sprintf("%v\t%v\t%v\tR%v, #%v", toPrint, progCount, tempArr2[0], tempArr2[2], tempArr2[1])
		}
	case 9:
		toPrint = fmt.Sprintf("%.9s %.2s %.16s %.5s", aString[0:9], aString[9:11], aString[11:27], aString[27:])
		tempArr1 := strings.Fields(toPrint)
		var tempArr2 [4]string

		for i := 0; i <= 15; i++ {
			if codeMatrix[i] == tempArr1[0] {
				tempArr2[0] = instMatrix[i]
			}
		}
		for i := 1; i < 4; i++ {
			decimal, err := binaryToDecimal(tempArr1[i])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				tempArr2[i] = strconv.FormatInt(decimal, 10)
			}
		}
		intValue, err := strconv.Atoi(tempArr2[1])
		if err != nil {
			// Handle the error if the conversion fails
			fmt.Println("Error:", err)
			//return
		}
		intValue = intValue * 16
		strValue := strconv.Itoa(intValue)

		tempArr2[1] = strValue

		toPrint = fmt.Sprintf("%v \t%v\t%v\tR%v, %v, LSL %v", toPrint, progCount, tempArr2[0], tempArr2[3], tempArr2[2], tempArr2[1])
	case 10:
		toPrint = fmt.Sprintf("%.10s %.12s %.5s %.5s", aString[0:10], aString[10:22], aString[22:27], aString[27:])
		tempArr1 := strings.Fields(toPrint) //stores each split from aString in an array
		var tempArr2 [4]string              //will store the instruction according to tempArr1

		for i := 0; i <= 10; i++ {
			if codeMatrix[i] == tempArr1[0] {
				tempArr2[0] = instMatrix[i]
			}
		}
		for i := 1; i < 4; i++ {
			decimal, err := binaryToDecimal(tempArr1[i])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				tempArr2[i] = strconv.FormatInt(decimal, 10)
			}
		}

		toPrint = fmt.Sprintf("%v  \t%v\t%v\tR%v, R%v, #%v", toPrint, progCount, tempArr2[0], tempArr2[3], tempArr2[2], tempArr2[1])
	case 11:
		if aString[0:5] == "11111" {
			toPrint = fmt.Sprintf("%.11s %.9s %.2s %.5s %.5s", aString[0:11], aString[11:20], aString[20:22], aString[22:27], aString[27:])
			tempArr1 := strings.Fields(toPrint) //stores each split from aString in an array
			var tempArr2 [5]string

			for i := 0; i <= 15; i++ {
				if codeMatrix[i] == tempArr1[0] {
					tempArr2[0] = instMatrix[i]
				}
			}
			for i := 1; i < 5; i++ {
				decimal, err := binaryToDecimal(tempArr1[i])
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				} else {
					tempArr2[i] = strconv.FormatInt(decimal, 10)
				}
			}
			toPrint = fmt.Sprintf("%v\t%v\t%v\tR%v, [R%v, #%v]", toPrint, progCount, tempArr2[0], tempArr2[4], tempArr2[3], tempArr2[1])
		} else {
			toPrint = fmt.Sprintf("%.11s %.5s %.6s %.5s %.5s", aString[0:11], aString[11:16], aString[16:22], aString[22:27], aString[27:])
			tempArr1 := strings.Fields(toPrint) //stores each split from aString in an array
			var tempArr2 [5]string              //will store the instruction according to tempArr1
			for i := 0; i <= 17; i++ {
				if codeMatrix[i] == tempArr1[0] {
					tempArr2[0] = instMatrix[i]
				}
			}
			for i := 1; i < 5; i++ {
				decimal, err := binaryToDecimal(tempArr1[i])
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				} else {
					tempArr2[i] = strconv.FormatInt(decimal, 10)
				}
			}
			if tempArr2[0] == "LSL" || tempArr2[0] == "LSR" || tempArr2[0] == "ASR" {
				toPrint = fmt.Sprintf("%v\t%v\t%v\tR%v, R%v, #%v", toPrint, progCount, tempArr2[0], tempArr2[4], tempArr2[3], tempArr2[2])
				//fmt.Println("lol")
			} else {
				toPrint = fmt.Sprintf("%v\t%v\t%v\tR%v, R%v, R%v", toPrint, progCount, tempArr2[0], tempArr2[4], tempArr2[3], tempArr2[1])
			}
		}
	case 12:
		toPrint = aString

		result, err := signedBinToDecimal(aString)
		if err != nil {
			fmt.Println(err)
		}
		toPrint = fmt.Sprintf("%v  %v\t%v", toPrint, progCount, result)
	default:
		toPrint = aString

		result, err := signedBinToDecimal(aString)
		if err != nil {
			fmt.Println(err)
		}
		toPrint = fmt.Sprintf("%v  %v\t%v", toPrint, progCount, result)
	}

	return toPrint
}

/*
binaryToDecimal: this function converts binary numbers to decimal
used for instructions that cannot be negative
*/
func binaryToDecimal(binary string) (int64, error) {
	decimal, err := strconv.ParseInt(binary, 2, 64)
	if err != nil {
		return 0, err
	}
	return decimal, nil
}

/*
signedBinToDecimal: this function converts signed binary number to decimal
*/
func signedBinToDecimal(binaryString string) (*big.Int, error) {
	isNegative := binaryString[0] == '1'

	if isNegative { // if it's negative, calculate its two's complement
		// invert all bits
		inverted := ""
		for _, bit := range binaryString {
			if bit == '0' {
				inverted += "1"
			} else {
				inverted += "0"
			}
		}

		// convert the inverted string to a big.Int
		bigInt := new(big.Int)
		_, success := bigInt.SetString(inverted, 2)
		if !success {
			return nil, fmt.Errorf("unknown binary string")
		}

		bigInt.Add(bigInt, big.NewInt(1)) // add 1 to get the two's complement

		bigInt.Neg(bigInt) // negate the result to make it negative

		return bigInt, nil
	}

	// if positive, convert the binary string directly to a big.Int
	bigInt := new(big.Int)
	_, success := bigInt.SetString(binaryString, 2)
	if !success { //error handling
		return nil, fmt.Errorf("unknown binary string")
	}
	return bigInt, nil
}

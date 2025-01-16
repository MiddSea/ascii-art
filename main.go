package main

// package main

import (
	// "bufio"
	"fmt"
	"log"

	// "slices"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	SPACE      = " "
	SNGL_QUOTE = "'"
)

func main() {
	err := checkArgs(os.Args)
	checkError(err)

	sampleInFile := os.Args[1]
	resultOutFile := os.Args[2]

	sampleText, err := readSampleFile(sampleInFile)
	checkError(err)

	resultText, err := processText(sampleText)
	checkError(err)

	err = writeResult(resultOutFile, resultText)
	checkError(err)
	fmt.Println("Conversion successful. Result saved in:", resultOutFile)

}

func checkArgs(args []string) (err error) {
	switch len(args) {
	case 1:
		err = fmt.Errorf("no input / output files given")
		return err // return false and error message if less than 2 arguments
	case 2:
		err = fmt.Errorf("only input, no output file given")
		return err // return false and error message if less than 3 arguments
	case 3:
		return nil // return true if 3 arguments
	default:
		err = fmt.Errorf("too many arguments")
		return err // return false and error message if more than 3 arguments
	}
}

func readSampleFile(filename string) (content string, err error) {
	var contentB []byte
	contentB, err = os.ReadFile(filename)
	if err != nil {
		log.Panicf("failed to read sample file: %v", err) // needs log package
	}
	return string(contentB), err

}

func writeResult(filename string, content string) error {
	err := os.WriteFile(filename, []byte(content), 0644) //
	if err != nil {
		log.Panicf("failed to writeResult to file: %v", err) // needs log package
	}
	return err
}

func checkError(err error) {
	if err != nil {
		log.Panicf("error: %v", err)
	}
}

// processText applies all transformations on the text
func processText(text string) (oText string, err error) {
	text, err = processHex(text)
	checkError(err)
	text, err = processBin(text)
	checkError(err)
	text, err = processCase(text)
	checkError(err)
	text, err = processPunctuation(text)
	checkError(err)
	text, err = processQuotesRegEx(text)
	checkError(err)
	text, err = processAtoAn(text)
	checkError(err)
	return text, err
}

func processHex(text string) (oText string, err error) {
	re := regexp.MustCompile(`(\w+)\s?\(hex\)`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		parts := strings.Fields(match)
		hexVal := parts[0]
		decVal, err := strconv.ParseInt(hexVal, 16, 0)
		if err != nil {
			return match
		}
		return fmt.Sprintf("%d", decVal)
	}), nil
}

// Convert binary numbers to decimal
func processBin(text string) (oText string, err error) {
	re := regexp.MustCompile(`(\w+)\s?\(bin\)`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		parts := strings.Fields(match)
		binVal := parts[0]
		decVal, err := strconv.ParseInt(binVal, 2, 0)
		if err != nil {
			return match // If there's an error, return the original
		}
		return fmt.Sprintf("%d", decVal)
	}), nil
}

func processCase(text string) (oText string, err error) {
	// Handle simple (up), (low), (cap)
	re := regexp.MustCompile(`(\w+)\s?\((up|low|cap)\)`)
	text = re.ReplaceAllStringFunc(text, func(match string) string {
		parts := strings.Fields(match)
		word := parts[0]
		command := parts[1]

		switch command {
		// was here round brakckets
		case "(up)":
			return strings.ToUpper(word)
		case "(low)":
			return strings.ToLower(word)
		case "(cap)":
			return capitalize(word)
		}
		return match
	})

	// Handle case transformations with a number (up, 2), (low, 3)
	reWithNum := regexp.MustCompile(`(\w+\s?)+\((up|low|cap),\s?(\d+)\)`)
	text = reWithNum.ReplaceAllStringFunc(text, func(match string) string {
		parts := strings.Fields(match)
		command := strings.Trim(parts[len(parts)-2], `\(,`)
		// numWords, _ := strconv.Atoi(parts[2])
		numWords, _ /*err*/ := strconv.Atoi(strings.Trim((parts[len(parts)-1]), `\)`))
		// Split the text into words and apply the transformation to the first 'numWords' words
		iLastWord := len(parts) - 2
		iFirstWord := iLastWord - numWords
		if iFirstWord < 0 {
			iFirstWord = 0 // Ensure the first word Idx is not negative
			// err = fmt.Errorf("number of words to transform is greater than the number of words in the text")
		}
		for i := iFirstWord; i <= iLastWord; i++ {
			switch command {
			case "up":
				parts[i] = strings.ToUpper(parts[i])
			case "low":
				parts[i] = strings.ToLower(parts[i])
			case "cap":
				parts[i] = capitalize(parts[i])
			}
		}
		return strings.Join(parts[:iLastWord], " ")
	})
	return text, err
}

func processPunctuation(text string) (oText string, err error) {
	// Remove space before punctuation

	// move space between punctuation to the right of the punctuation
	text = regexp.MustCompile(`([.,!?;]+)\s+([.,!?;]+)`).ReplaceAllString(text, "${1}${2} ")

	// move space left of punctuation to the right of the punctuation
	text = regexp.MustCompile(`\s+([.,!?;]+)`).ReplaceAllString(text, "${1} ")

	// cut multiple spaces in text and trim ends
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	return text, nil
}

func inRange(rng []rune, i int) bool {
	return i >= 0 && i <= len(rng)
}

func between(rng []rune, lower, higher int) bool {
	return lower <= higher && inRange(rng, lower) && inRange(rng, higher)
}

func isSpace(txt []rune, i int) bool {
	return txt[i] == ' '
}

func isQuote(txt []rune, i int) bool {
	return txt[i] == '\''
}

func isSpQuote(txt []rune, i int) bool {
	return isQuote(txt, i) || isSpace(txt, i)
}

func printRuler() {
	fmt.Println("|0         1         2         3         4         5         6         7         8")
	fmt.Println("|012345678901234567890123456789012345678901234567890123456789012345678901234567890")
}

// processQuotesRegEx processes single quotes according to the specified rules
func processQuotesRegEx(text string) (string, error) {
	// Regular expression to match single quotes with spaces around them and
	// and beginning and end of text
	// TO POSSIBLY ADD: dealing with multiple quotes in a row
	// allQuotesRegEx := regexp.MustCompile(`(^\s*'\s*)|(\s+'\s*)|(\s*'\s+)|(\s*'\s*$)`)
	// allQuotesRegEx := regexp.MustCompile(`(^\s*'\s*)|('\s+')|(\s+'\s*)|(\s*'\s+)|(\s*'\s*$)`)
	// 2025-01-13 17-10 double quotes removed           ^^^^^^^
	allQuotesRegEx := regexp.MustCompile(`(^\s*'\s*)|(\s+'\s*)|(\s*'\s+)|(\s*'\s*$)`)
	// qtSpQtRegEx := regexp.MustCompile(`'\s+'`)
	bgIdxStrRegEx := regexp.MustCompile(`^.`)
	enIdxStrRegEx := regexp.MustCompile(`.$`)
	bgIdxStr := bgIdxStrRegEx.FindStringIndex(text)
	enIdxStr := enIdxStrRegEx.FindStringIndex(text)
	fstTxTChar := text[bgIdxStr[0]:bgIdxStr[1]]
	lstTxTChar := text[enIdxStr[0]:enIdxStr[1]]
	var testResult strings.Builder
	var testResultStr string
	// var result strings.Builder
	//redOnWhite := "\033[91m\033[107m" //red text on White background
	// rOw := "\033[91m\033[107m" //red text on White background
	// resetColors := "\033[0m" // reset text output colours to default
	// lightGray := "\033[37m" // light gray text
	rOlg := "\033[91m\033[47m" // red text on light gray background
	// lightGrayBackground := "\033[47m" // light gray background
	dfCol := "\033[0m" // default text color

	// fmt.Printf ("\033[91m\033[107mThis is red text on a white background\033[0m\n")
	// fmt.Printf (redOnWhite+"This is more red text on a white background"+resetColors+"\n")
	fmt.Printf("bgIdxStr: %02v fstTxTChar %s%s%s enIdxStr: %-2v lstTxTChar %s%s%s \n",
		bgIdxStr, rOlg, fstTxTChar, dfCol, enIdxStr, rOlg, lstTxTChar, dfCol)
	if enIdxStr == nil {
		return "", fmt.Errorf("empty text")
	}

	qtMatches := allQuotesRegEx.FindAllStringIndex(text, -1)
	qtMatchesCount := len(qtMatches)
	printRuler()
	fmt.Printf("|%s|\n", text)

	// trial run quick and dirty
	// testResult := ""
	bgTstOrgStr := -1
	enTstOrgStr := -1
	qt_col := fmt.Sprintf("%s'%s", rOlg, dfCol)
	sp_col := fmt.Sprintf("%s %s", rOlg, dfCol)
	// spQtQtSp_col := fmt.Sprintf("%s '' %s", rOlg, dfCol)
	for i, qtMatch := range qtMatches {
		fmt.Printf("i: %d, qtMatch: %v\n", i, qtMatch)
		// /////////////
		// first qt Match
		if i == 0 { // first quoteMatch
			if qtMatch[0] == 0 { // first quote opening quote at beginning of text
				testResult.WriteString("'")
				testResultStr += qt_col
				fmt.Printf("X%sX", testResultStr)
				bgTstOrgStr = qtMatch[1] // set beginning of original string to end of quote match
			} else if qtMatch[0] > 0 { // text at beginning then first qtMatch 
				// bgTstOrgStr = 0
				enTstOrgStr = qtMatch[0]
				testResult.WriteString(text[:enTstOrgStr])
				testResult.WriteString(SPACE+SNGL_QUOTE)
				testResultStr += text[:enTstOrgStr] + sp_col + qt_col
				fmt.Printf("T%sT", testResultStr)
				bgTstOrgStr = qtMatch[1] // set beginning of original string
	
				// TO DO
			}
		// middle qt Match
		} else if i > 0 && i < qtMatchesCount { // middle qt Match
			if i < qtMatchesCount && i%2 == 0 && i > 0 { // even qtMatch = closing quote
				testResult.WriteString("E")
				testResult.WriteString(text[bgTstOrgStr:qtMatch[0]])
				testResultStr += "E"
				testResultStr += text[bgTstOrgStr:qtMatch[0]]	
		   }
		// last qt Match
		} else if i == qtMatchesCount { // last qt Match
		
		}
        // ////////////////////////////////

//// TO DO DO DO
		    else if i == qtMatchesCount { // last quote
			// last quote
			fmt.Printf("L")
			testResult.WriteString("L")
			testResult.WriteString(text[bgTstOrgStr:qtMatch[0]])
			testResult.WriteString(SNGL_QUOTE)
			testResultStr += "L"
			testResultStr += text[bgTstOrgStr:qtMatch[0]] + qt_col
			fmt.Printf("L%sL", testResultStr)	
		} else if i%2 != 0 { // odd qtMatch
		    testResult.WriteString("°")
			testResult.WriteString(text[bgTstOrgStr:qtMatch[0]])
			testResultStr += "°"
			testResultStr += text[bgTstOrgStr:qtMatch[0]]
			// testResult.WriteString("°")
			fmt.Printf("O%sO", testResultStr)
		} // else if 

		fmt.Printf(">i:%d %02v %02v |%s|<\n", i, qtMatches, qtMatch, text[qtMatch[0]:qtMatch[1]])
		fmt.Printf("|%s|\n", testResultStr)
	}
	
	return testResult.String(), nil
	// fmt.Printf("qtMatchesCount: %v, qtMatches: %v\n",
	//	qtMatchesCount, qtMatches)
	//
	//if qtMatchesCount == 0 {
	//	fmt.Printf(" *** no quotes in text ***\n")
	//	return text, nil // do nothing if no quotes
	//}
	//if qtMatchesCount%2 != 0 {
	//	fmt.Printf(" *** odd number of valid quotes ***\n")
	//	return "", nil /* fmt.Errorf("odd number of valid quotes, didn't count apostrophes in text")*/
	//}
	//var bgQtIdx []int = qtMatches[0]
	//var enQtIdx []int = qtMatches[qtMatchesCount-1]
	//var fstChar string = text[bgQtIdx[1]-1 : bgQtIdx[1]]
	// // var fstCharIdx int = bgQtIdx[1] // DEBUG actually [0]
	// var fstCharIdx int = bgQtIdx[0] // DEBUG
	//var lstChar string = text[enQtIdx[1]-1 : enQtIdx[1]]
	// var qtTxtIdx int = strings.Index(text, "'") // DEBUG
	// Idx of last qtMatched character
	//fmt.Printf("< enQtIdx %v |%s|, bgQtIdx: %v, |%s|, \n", // DEBUG
	//enQtIdx, lstChar, bgQtIdx, fstChar) // DEBUG
	// fmt.Printf("bgQtIdx: %v, enQtIdx: %v qtTxtIdx\n", bgQtIdx, fstChar, enQtIdx, qtTxtIdx, lstChar)
	// enIdxStr           // Idx of last character in text
	// se
	// lastChar = text[:finalTextIdx]                                 // DEBUG
	//return result.String(), nil
	//
	//bgOrgStr := -1   // original string begin Idx
	//enOrgStr := -1   // original string end Idx
	//inQuote := false // correct until PAST first there is a qtMatch
	//
	//for qtMatchNo, qtMatch := range qtMatches {
	//	// var bgQtMatch, enQtMatch
	//	bgQtMatch, enQtMatch := qtMatch[0], qtMatch[1]
	//	// TO DO uses simple integers instead of slices
	//
	//	if qtMatchNo == 0 { // first qtMatch
	//
	//		if bgQtMatch > 0 { // text before first qtMatch
	//			bgOrgStr = 0
	//			enOrgStr = bgQtMatch
	//			// copy text before first qtMatch
	//			result.WriteString(text[bgOrgStr:enOrgStr])
	//			fmt.Printf(">>>> qt Match > 0 result: |%s| ", result.String()) // DEBUG
	//			// orevious text end is the beginning of the first qtMatch
	//			bgOrgStr = enQtMatch
	//			inQuote = true
	//		} else { //quote at beginning
	//			// result.WriteString("'") // write quote
	//			bgOrgStr = enQtMatch
	//			inQuote = true
	//		}
	//		fmt.Printf("bgQtMatch: %v, enQtMatch: %v \n",
	//			bgQtMatch, enQtMatch)
	//	} else if qtMatchNo <= qtMatchesCount-2 { // not the last qtMatch
	//		// do stuff for middle qtMatches
	//		// startNext, endNext := matches[matchNo+1][0], matches[matchNo+1][1]
	//		fmt.Printf("bgQtMatch: %v, enQtMatch: %v \n",
	//			bgQtMatch, enQtMatch)
	//		fmt.Printf(">#%02d %2v!!!",
	//			qtMatchNo, qtMatches[qtMatchNo])
	//
	//		if inQuote { // Opening quote
	//			fmt.Printf(" if inQuote >")
	//			if !(bgOrgStr == 1 && qtMatchNo == 1) /*&& text[bgQtMatch-1] != ' '*/ {
	//				result.WriteString("#")
	//			}
	//			result.WriteString("°")
	//			fmt.Printf("text[bgQtMatch:enQtMatch]: |%s|\n", text[bgQtMatch:enQtMatch])
	//			result.WriteString(text[bgOrgStr:bgQtMatch])
	//			inQuote = false
	//		} else { // Closing quote
	//			fmt.Printf(" else !inQuote")
	//			//  result.WriteString(text[bgOrgStr:bgQtMatch])
	//			result.WriteString(">^") // ^ // was here
	//			if enQtMatch < len(text) && text[enQtMatch] != ' ' {
	//				result.WriteString("_")
	//			}
	//			result.WriteString(text[bgOrgStr:bgQtMatch])
	//			inQuote = true
	//
	//		}
	//
	//	} else /* shouldn't ever get here.*/ { // last qtMatch
	//		// do stuff for last qtMatch
	//		// copy text after last qtMatch
	//		// enOrgStr = lstTxTChar
	//
	//		bgOrgStr = enQtMatch // TO DO check if this is needed
	//		result.WriteString("&&&&&&")
	//		result.WriteString(text[bgOrgStr:enOrgStr])
	//	}
	//	bgOrgStr = enQtMatch // moved this up
	//	// fmt.Printf(" bgOrgStr: %v, enQtMatch:%v, qtMatchNo %v, text[bgOrgStr:enOrgStr]: %v\n", // DEBUG bounds check
	//	// bgOrgStr, enOrgStr, qtMatchNo, text[bgOrgStr:enOrgStr]) // DEBUG bounds check
	//}
	//	result.WriteString(text[bgOrgStr:]) // TO DO check if this is needed
	//	return result.String(), nil
}

// processAtoAn replaces 'a' or 'A' followed by a word starting with a vowel or a 'h/H' with 'an' or 'An'
func processAtoAn(text string) (oText string, err error) {
	re := regexp.MustCompile(`\b(a|A)\s+([aeiouhAEIOUH]\w*)\b`)
	return re.ReplaceAllString(text, "${1}n $2"), nil
}

// Capitalize the first letter of a word
func capitalize(word string) string {
	if len(word) == 0 {
		return word
	}
	return strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
}

package main
// package main

import (
	// "bufio"
	"fmt"
    "log"
	// "slices"
	"os"
	"regexp"
	"strings"
	"strconv"
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

/* old main function without proper file handling
func main() {
    examples := []string{        
		"This is so exciting (up, 2) and naughty",
        "1E (hex) files were added",
        "It has been 10 (bin) years",
        "Ready, set, go (up)!",
		"A offer, A Offence and a hotel",
        "I should stop SHOUTING (low)",
        "Welcome to the Brooklyn bridge (cap, 7)",
    }

	if err := processFile("input.txt", "output.txt"); err != nil {
		fmt.Println("Error processing file:", err)
	}

	for _, example := range examples {
		fmt.Println("Original:", example)
		fmt.Println("Processed:", processText(example))
		fmt.Println()
	}
}
*/


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
		log.Panicf("failed to readSample  file: %v", err) // needs log package
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
/*
// processFile reads the input file, processes each line, and writes the output
func processFile(inputFile, outputFile string) error {
	inFile, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		processedText := processText(scanner.Text())
		_, err := outFile.WriteString(processedText + "\n")
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}*/


// processText applies all transformations on the text
func processText(text string) (oText string, err error) {
	text, err = processHex(text) 			; checkError(err)
	text, err = processBin(text) 			; checkError(err)
	text, err = processCase(text) 			; checkError(err)
	text, err = processPunctuation(text) 	; checkError(err)
	text, err = processAtoAn(text) 			; checkError(err)
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
	 // reWithNum := regexp.MustCompile(`(\w+)\s?\((up|low|cap),\s?(\d+)\)`)
 	//reWithNum := regexp.MustCompile(`(\w+\s?)+\((up|low|cap),\s?(\d+)\)`)
 	//text = reWithNum.ReplaceAllStringFunc(text, func(match string) string {
	// parts := strings.Fields(match)
	// word := parts[0]
	// command := parts[1]
	// numWords, _ := strconv.Atoi(parts[2])

	 reWithNum := regexp.MustCompile(`(\w+\s?)+\((up|low|cap),\s?(\d+)\)`)
	 text = reWithNum.ReplaceAllStringFunc(text, func(match string) string {
		 parts := strings.Fields(match)
		 // word := parts[0]
		 
		 // command := parts[1]
		 command := strings.Trim(parts[len(parts)-2], `\(,`)
		 // numWords, _ := strconv.Atoi(parts[2])
		 numWords, _ := strconv.Atoi(strings.Trim((parts[len(parts)-1]), `\)`))
	  // Split the text into words and apply the transformation to the first 'numWords' words
	     iLastWord := len(parts) - 2
		 iFirstWord := iLastWord - numWords
		 if iFirstWord < 0 {
			 iFirstWord = 0
		 }	
		// words := parts[iFirstWord:iLastWord]
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
	// text = regexp.MustCompile(`\s([.,!?;])`).ReplaceAllString(text, "${1}")
	// Remove space between punctuation
	// text = regexp.MustCompile(`\s+([.,!?;])+\s*?([.,!?;])+`).ReplaceAllString(text, "${1}${2}")
	
	// move space between punctuation to the right of the punctuation
	text = regexp.MustCompile(`([.,!?;]+)\s+([.,!?;]+)`).ReplaceAllString(text, "${1}${2} ")
	
	// move space left of punctuation to the right of the punctuation
	text = regexp.MustCompile(`\s+([.,!?;]+)`).ReplaceAllString(text, "${1} ")

	// cut multiple spaces in text and trim ends
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	
	// Handle punctuation marks followed by a space
//	text = regexp.MustCompile(`([.,!?;]+)\s+`).ReplaceAllString(text, "S${1}T ")
//	text = regexp.MustCompile(`([.,!?;]+)([^.,!?:;\s])`).ReplaceAllString(text, "U${1} ${2}V")

	// Handle special cases like ellipses
//	text = regexp.MustCompile(`\.\.\.`).ReplaceAllString(text, "...")
//	text = regexp.MustCompile(`\!\?\?`).ReplaceAllString(text, "!?")
	return text, nil
}

// processQuotes removes extra spaces around quotes
// opening single quotes with a space before or beginning the string 
// and closing single quotes with a space after, unless at the end of the string 
/* func processQuotes(text string) (oText string, err error) {
	// Handle single quotes
	quoteWithSpace = regexp.MustCompile(`(')\s+`)
}*/
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


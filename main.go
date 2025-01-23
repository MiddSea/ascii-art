package main

// package main

import (
	// "bufio"
	"fmt"
	"log"

	// "slices"
	"os"
	//"regexp"
	// "strconv"
	"strings"
)

var err error

const (
	SPACE      = " "
	SNGL_QUOTE = "'"
	USAGE	   = "USAGE: go run . \"An ASCII string, with 123 \\n and Symbols +0!#*\""
	BANNERFILE = "standard.txt"
)
	
func main() {
		// 1. Check and validate command line arguments
		if err := checkArgs(os.Args); err != nil {
			fmt.Printf("%v\n", err)
			fmt.Println(USAGE)
		}
		checkError(err)

		// 2. Get input string from arguments
		asciiArtInput := os.Args[1]

		// 3. Read banner file and parse character templates
		// bannerFile := "standard.txt"  // default banner defined above
		bannerData, err := readAsciiFormatData(BANNERFILE)
		if err != nil {
			fmt.Printf("Error reading banner file: %v\n", err)
		}
		checkError(err)

		// 4. Generate ASCII art
		// TO DO: Implement generateAsciiArt function
		result, err := generateAsciiArt(asciiArtInput, bannerData)
		if err != nil {
			fmt.Printf("Error generating ASCII art: %v\n", err)
		}
		checkError(err)

		// 5. Print result
		fmt.Print(result)
	}

func checkArgs(args []string) (err error) {
	switch len(args) {
	case 1:
		err = fmt.Errorf("no ASCII string for banner provided")
		return err // return false and error message if less than 2 arguments
	case 2:
		asciiArtInput := args[1]
		for _, r := range asciiArtInput {
			if (r < 32 && r != 10) || r > 127 {
				if r < 32 && r != 10 {
					err = fmt.Errorf("string contains non-printable characters")
				}
				if r > 127 {
					if err != nil {
						err = fmt.Errorf("%v and non-ASCII characters", err)
					} else {
						err = fmt.Errorf("string contains non-ASCII characters")
					}
				}
			}
		}
		return err // return false and error message if non-printable characters or non-ASCII characters
	default:
		err = fmt.Errorf("too many arguments")
		return err // return false and error message if more than 3 arguments
	}
}

func readAsciiFormatData(filename string) (asciiFormatStr string, err error) {
	var asciiFormatData []byte
	asciiFormatData, err = os.ReadFile(filename)
	if err != nil {
		log.Panicf("failed to read sample file: %v", err) // needs log package
	}
	return string(asciiFormatData), err

}


func checkError(err error) {
	if err != nil {
		log.Panicf("error: %v", err)
	}
}

func generateAsciiArt(bannerInput string, asciiFmtData string) (banner string, err error) {
	// 1. Split input string into bannerInputLines
	bannerInputLines := strings.Split(bannerInput, "\n")
	//asciiChars := make([]string, 0)
	// asciiChars = strings.SplitAfter(asciiArtInput, "\n\n")
	// asciiFmtChars := strings.Split(bannerData, "\n\n")
	asciiFmtChars := strings.Split(asciiFmtData, "\n")
	// map the banner data to the ascii characters in the format [ascii character rune value] = [row 0-] banner character
	

	// 2. Initialize result
	banner = ""
	// 3. Loop through each line
	for _, bannerInputLine := range bannerInputLines {
		// 4. Loop through each character in the line
		for _, r := range bannerInputLine {
			// 5. Check if character is printable
			switch  {
			case r == '\n':

			case r >= 32 && r <= 127:
			    
				// 6. Get the index of the character in the banner data
				index := int(r - 32)
				// 7. Get the ASCII art for the character
				charArt := ""
				for i := 0; i < 8; i++ { 
					charArt = asciiFmtChars[index+1]
				}
				
				// 8. Append the ASCII art to the result
				// result = appen	//"regexp"d(result, charArt)
				banner = banner + charArt
			}
		}
		// 9. Append newline character to the result
		// result += "\n"
	}
	return banner, err


}


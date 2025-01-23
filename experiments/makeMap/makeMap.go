package main

import (
    "bufio"
    "fmt"
    "os"
)

func makeAsciiMap(filename string) (map[rune][]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    asciiMap := make(map[rune][]string)
    scanner := bufio.NewScanner(file)
    
    currentChar := rune(32) // Start with space character
    var lines []string

    for scanner.Scan() {
        line := scanner.Text()
        lines = append(lines, line)
        
        if len(lines) == 8 {
            asciiMap[currentChar] = make([]string, 8)
            copy(asciiMap[currentChar], lines)
            lines = lines[:0]
            currentChar++
        }
    }

    return asciiMap, scanner.Err()
}

func main() {
    fmt.Println("Basic map:")
    aBasicMap()
    fmt.Println("..done")
    asciiMap, err := makeAsciiMap("../standard.txt")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // Test print for '!' character
    if art, ok := asciiMap['!']; ok {
        fmt.Println("ASCII art for '!':")
        for _, line := range art {
            fmt.Printf("%s\n", line)
        }
    }
}

func aBasicMap() {
    asciiMapEx := make(map[rune][]string)
    asciiMapEx['!'] = []string{
        "   ",
        " ! ",
        " ! ",
        " ! ",
        "   ",
        " ! ",
        "   ",
        "   ",
    }
    for i := 0; i < 8; i++ {
        fmt.Printf("%s\n", asciiMapEx['!'][i])
    }
}
package text

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// TextSearcher structure
// TODO: add in objects/data structures required by the Search function
type TextSearcher struct {
	Text       []string
	Dictionary map[string][]int
}

// We need to check if we are in a test, if so go up one more directory level
func IsRunningTest() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}
	return false
}

// NewSearcher returns a TextSearcher from the given file
// TODO: Load/process the file so that the Search function work
func NewSearcher(filePath string) (*TextSearcher, error) {
	textSearcher := TextSearcher{
		Text:       []string{},
		Dictionary: make(map[string][]int),
	}
	var file *os.File
	var err error

	// tests are pathing to ~/PROGRAMMING-CHALLENGE/go rather than to ~/PROGRAMMING-CHALLENGE
	// tests should be pathed ../../ rather than ../ (tests are run based on test's directory)
	if IsRunningTest() {
		file, err = os.Open(fmt.Sprintf("../%s", filePath))
	} else {
		file, err = os.Open(filePath)
	}

	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	index := 0
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Println(line)
		// string split the line on spaces (to preserve punctuatation/formatting)
		splitLine := strings.Split(line, " ")
		for _, word := range splitLine {
			// There is the edge case for empty strings
			cleanedWord := sanitizeString(word)
			if strings.TrimSpace(cleanedWord) != "" {
				// we are adding the new index to the map
				textSearcher.Dictionary[cleanedWord] = append(textSearcher.Dictionary[cleanedWord], index)
				// adding unaltered word to the []string
				textSearcher.Text = append(textSearcher.Text, word)
				index++
			}
		}
	}

	// fmt.Println(textSearcher.Dictionary)
	// fmt.Println(textSearcher.Text)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return &textSearcher, nil
}

// before adding the string to the Dictionary, we want to sanitize the word to ensure that all permutations are
// stored in lowercase and are alphanumeric + dash/apostrophe per documentation.
func sanitizeString(word string) string {
	word = strings.ToLower(word)
	var result strings.Builder
	for _, char := range word {
		// check if the character is a-z, 0-9, apostrophe, and dash
		// (we can ignore capitals since we want to search case agnostic)
		if unicode.IsLower(char) || unicode.IsNumber(char) || char == '-' || char == '\'' {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// Search searchs the file loaded when NewSearcher was called for the given word and
// returns a list of matches surrounded by the given number of context words
// Solution 1:
// time_complexity: O(n) where n is the length of text document
// space_complexity: O(m) where m is the number of times the word appears in the text
// create a fifo queue (slice) in search bar for context, we will traverse the file exactly once
// we will populate the the queue with words until we hit the word we are searching for
// we will then create a new slice with the word at the front + context after the word
// we will then concatenate the slices together and return them to the user. (a sliding window)
// search would be o(N) since we would have to traverse the file every time

// Solution 2:
// time_complexity: O(1), we are searching map for the word (O(1) time), then gathering context around the word after
// space_complexity: O(n), we will need to create a map of the O(n) plus storing the file in memory O(n) totaling O(2n) == O(n)
// We could hold the entire file in memory based on the readme. Therefore, we can make a map[string][]int
// each time we see a word, we will add the index of the word to the map. ex: map["test"] = {3, 6, 7}
// We will hold another string array in memory for the text file itself (both of these structs will be created
// in the initializing funct). Therefore, once we run the searching algo, we can do a quick loopup of the word
// since we will know where the word is based on the index in the map

// TODO: Implement this function to pass the tests in search_test.go
func (ts *TextSearcher) Search(word string, context int) []string {
	search := []string{}
	// Search the map for the sanatized string, then create a window around it:
	//     map(index-context:index+context)
	// We need to ensure that there is no range out of bounds error
	// fmt.Println(ts.Dictionary[sanitizeString(word)])
	for _, index := range ts.Dictionary[sanitizeString(word)] {
		leftWindow := index - context
		if leftWindow < 0 {
			leftWindow = 0
		}
		rightWindow := index + context + 1
		if rightWindow > len(ts.Text) {
			rightWindow = len(ts.Text)
		}
		// fmt.Println(index)
		// fmt.Println(ts.Text[leftWindow:rightWindow])
		search = append(search, strings.Join(ts.Text[leftWindow:rightWindow], " "))
	}

	return search
}

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

/*

BRAHE VERSION 1.0
_________________
[0]. Model
[1]. Clear Screen Fn
[2]. Display HomeView Fn
[3]. Storage Setup Fn
[4]. Operation FNs
	- [4.1]. Write
	- [4.2]. Read
	- [4.3]. Update <PENDING>
	- [4.4]. Delete
*/

// model
type Word struct {
	Word       string
	Definition string
	Examples   []string
	Nouns      []string
}

// Clear Screen
func clearSC() {
	if runtime.GOOS == "darwin" {
		screen := exec.Command("clear")
		screen.Stdout = os.Stdout
		screen.Run()
	} else {
		screen := exec.Command("cls")
		screen.Stdout = os.Stdout
		screen.Run()
	}
}

// Home View
func Homeview() {
	clearSC()

	// heading
	var version string = "1.0"
	var title string = "Hello I'M BRAHE " + version
	var line string = "---------------------------------------------"
	fmt.Println(line)
	fmt.Println("           ", title)
	fmt.Println(line)

	// options
	fmt.Println("OPERATIONS:")
	options := []string{"write", "read", "update", "delete"}
	for _, value := range options {
		fmt.Println("- ", value)
	}
	fmt.Println(line)
}

// Storage Setup/Config
func checkDirExistence(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	// default return
	return false, err
}

func processStorage() string {
	_dir := "./store"

	// check if directory exists
	dirExists, err := checkDirExistence(_dir)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		if dirExists {
			return "./store - directory is ready for use.."
		} else if !dirExists {
			// creating ./store directory
			os.Mkdir(_dir, os.ModePerm)
			return "./store - directory is created and ready for use.."
		}
	}

	// default return
	return "Unable to verify storage!"
}

/*
*******************
OPERATION FNs
*******************
*/

// write
func write(word string, definition string, examples []string, nouns []string) {
	_word := Word{
		Word:       strings.TrimSuffix(word, "\n"),
		Definition: strings.TrimSuffix(definition, "\n"),
		Examples:   examples,
		Nouns:      nouns,
	}

	// converting struct into JSON
	encodedJSON, err := json.Marshal(_word)
	if err != nil {
		fmt.Println(err)
	} else {
		os.WriteFile("./store/"+strings.TrimSuffix(word, "\n")+".json", encodedJSON, os.ModePerm)
	}
}

// read
func read(word string) {
	// reading file from fs
	result, err := os.ReadFile("./store/" + word + ".json")

	// parsed json data
	var parsedWord Word
	if err != nil {
		fmt.Println(err)
	} else {
		stringifiedResult := string(result)
		err := json.Unmarshal([]byte(stringifiedResult), &parsedWord)
		if err != nil {
			fmt.Println(err)
		} else {
			// printing output
			Homeview()
			fmt.Println("Showing Result for", parsedWord.Word)
			fmt.Println("WORD       :", parsedWord.Word)
			fmt.Println("DEFINITION :", parsedWord.Definition)
			for index, value := range parsedWord.Examples {
				fmt.Println("EXAMPLE", index, " : "+value)
			}
			for index, value := range parsedWord.Nouns {
				fmt.Println("NOUN   ", index, " : "+value)
			}
			fmt.Println("---------------------------------------------")
		}
	}
}

// delete
func delete(word string) string {
	// checking if the file is available
	err := os.Remove("./store/" + word + ".json")
	if err != nil {
		fmt.Println(err)
		return "Error Occured While Performing Delete Operation.."
	}

	var alertMsg string = word + " is deleted from the ./store directory.."
	return alertMsg
}

func main() {
	sc := bufio.NewReader(os.Stdin)

	// setup storage folder
	storageStatus := processStorage()

	// display HomeView
	Homeview()
	fmt.Println(storageStatus)

	// allowing users to perform operations
	for {
		// Homeview()
		fmt.Println("ENTER AN OPERATION:")

		// Operation input
		var option string
		fmt.Scanln(&option)

		switch option {
		case "write":
			// taking user input for write operation
			for {
				Homeview()

				var (
					examples []string
					nouns    []string
				)

				// word input
				fmt.Println("Enter the word:")
				word, err := sc.ReadString('\n')
				if err != nil {
					fmt.Println(err)
				}
				if word == "" {
					break
				}

				// definition input
				fmt.Println("Enter definition:")
				definition, err := sc.ReadString('\n')
				if err != nil {
					fmt.Println(err)
				}
				if definition == "" {
					break
				}

				// examples input
				for {
					fmt.Println("Enter example:")
					example, err := sc.ReadString('\n')
					example = strings.TrimSpace(example)
					if err != nil {
						fmt.Println(err)
					}
					if example == "" {
						break
					}
					examples = append(examples, example)
				}

				// nouns input
				for {
					fmt.Println("Enter a noun:")
					noun, err := sc.ReadString('\n')
					noun = strings.TrimSpace(noun)
					if err != nil {
						fmt.Println(err)
					}
					if noun == "" {
						break
					}
					nouns = append(nouns, noun)
				}

				// writing to file-system
				write(word, definition, examples, nouns)
				break
			}
		case "read":
			// taking user input for read operation
			for {
				Homeview()
				fmt.Println("Enter the word to search:")
				readInput, err := sc.ReadString('\n')
				if err != nil {
					fmt.Println(err)
				}

				// running read operation
				read(strings.TrimSuffix(readInput, "\n"))
				break
			}
		case "update":
			fmt.Println("updating")
		case "delete":
			for {
				Homeview()
				fmt.Println("Enter the word to delete:")
				readInput, err := sc.ReadString('\n')
				result := delete(strings.TrimSuffix(readInput, "\n"))
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(result)
				}
				break
			}
		default:
			Homeview()
			fmt.Println("ALERT: '"+option+"'", " was an invalid input!")
			continue
		}
	}
}

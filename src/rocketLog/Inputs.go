package main

import (
	"bufio"
	"os"
	"log"
	"encoding/json"
	"path/filepath"
)

const STATE_FILE = "./state.json"

type FileInput struct {
	scanner     bufio.Scanner
	line_number int
	abs_path    string
	file *os.File
}


func NewFileInput(path string) *FileInput {
	file, err := os.Open(path)
	if(err != nil){
		log.Fatal(err)
	}

	abs_path, err := filepath.Abs(path)
	if(err != nil){
		log.Fatal(err)
	}

	file_scanner := bufio.NewScanner(file)
	fin := &FileInput{
		scanner: *file_scanner,
		abs_path: abs_path,
		file: file,
	}

	file_state := fin.loadState()

	fin.SkipTo(file_state[fin.abs_path])

	return fin
}

func (input *FileInput) SkipTo(skip_to int){
	for input.line_number < skip_to {
		log.Print("Current Line: ", input.ReadLine(), ", Input.line_number, ", input.line_number, " skip_to, ", skip_to)
	}
}

func (input *FileInput) Close() {
	input.file.Close()
	input.saveState()
}

func (input *FileInput) HasLine() bool {
	return input.scanner.Scan()
}

func (input *FileInput) ReadLine() string {
	if (input.scanner.Scan() == false){
		log.Fatal("No Tokens Left")
	}
	input.line_number++
	return input.scanner.Text()
}

func (input *FileInput) saveState(){
	file_map := input.loadState()
	file_map[input.abs_path] = input.line_number

	file, err := os.OpenFile(STATE_FILE, os.O_RDWR | os.O_TRUNC, 0666)
	if(err != nil){
		log.Fatal(err)
	}

	e := json.NewEncoder(file)
	err = e.Encode(file_map)
	if(err != nil){
		log.Fatal(err)
	}

	file.Close()
}

func (input *FileInput) loadState() map[string] int{
	var file *os.File
	var err error

	file, err = os.OpenFile(STATE_FILE, os.O_CREATE | os.O_RDWR, 0666)
	if(err != nil){
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	file_state := make(map[string]int)
	decoder.Decode(&file_state)
	return file_state
}


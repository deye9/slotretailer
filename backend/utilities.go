package service

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

// BasePath returns the base working path without the trailing path separator
func BasePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Error reading the BasePath. ", err)
	}

	// Split by the build path as we are not saving into the root
	return strings.Split(cwd, "\\build")[0]
}

// ReadFile reads a file and returns the data or an error
func ReadFile(filepath string) ([]byte, error) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		log.Fatalln("Reading File "+filepath+" failed.", err)
	}
	return ioutil.ReadFile(filepath)
}

// CheckError catches and handles backend errors
func CheckError(message string, err error, fatal bool) {
	if err != nil && fatal == false {
		go func() {
			// Write the error to the File System.
			WriteFile(BasePath()+"/build/error.log", []byte(time.Now().String()+" "+message+": "+err.Error()+"\n"))
		}()
		return
	} else if err != nil && fatal == true {
		go func() {
			// Write the error to the File System.
			WriteFile(BasePath()+"/build/error.log", []byte(time.Now().String()+" "+message+": "+err.Error()+"\n"))
			os.Exit(1)
		}()
	}

	// catch to error.
	return
}

// WriteFile writes the data into the file specified
func WriteFile(filepath string, data []byte) (bool, error) {

	// If the file doesn't exist, create it, or append to the file
	f, _ := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if _, err := f.Write(data); err != nil {
		f.Close() // ignore error; Write error takes precedence
	}

	if err := f.Close(); err != nil {
		// CheckErr("Close File Error.", err, 4)
	}

	return true, nil
}

// loads an html file for the user
func loadFile(fileName string) string {
	fileName = BasePath() + fileName

	if content, err := ReadFile(fileName); err == nil {
		return string(content)
	}

	return ""
}
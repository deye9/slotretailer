package service

import (
	"fmt"
	"io/ioutil"
	"log"
	smtp "net/smtp"
	"os"
	"strings"
	"time"
)

// BasePath returns the base working path without the trailing path separator
func BasePath() string {
	cwd, err := os.Executable()
	if err != nil {
		log.Fatalln("Error reading the BasePath. ", err)
	}

	// Split by the build path as we are not saving into the root {window: \\, mac: /}
	return strings.Split(cwd, "/build")[0]
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
	if err != nil && fatal {
		// Write the error to the File System.
		go WriteFile(BasePath()+"/error.log", []byte("Fatal Error: "+time.Now().String()+" "+message+": "+err.Error()+"\n"))
		os.Exit(1)
	}

	if err != nil {
		// Write the error to the File System.
		WriteFile(BasePath()+"/error.log", []byte(time.Now().String()+" "+message+": "+err.Error()+"\n"))
		return
	}
}

// WriteFile writes the data into the file specified
func WriteFile(filepath string, data []byte) (bool, error) {
	// If the file doesn't exist, create it, or append to the file
	f, _ := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if _, err := f.Write(data); err != nil {
		f.Close() // ignore error; Write error takes precedence
	}

	if err := f.Close(); err != nil {
		fmt.Println("Error Closing file: "+filepath, err)
		// CheckErr("Close File Error.", err, 4)
	}

	return true, nil
}

// WeekRange returns the start and end date of the selected week
func WeekRange(year, week int) (start, end time.Time) {
	start = WeekStart(year, week)
	end = start.AddDate(0, 0, 6)
	return
}

// WeekStart returns the start date of the week
func WeekStart(year, week int) time.Time {
	// Start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	// Roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// Difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}

func lastPeriod(t time.Time, period time.Month) (start, end time.Time) {
	y, m, _ := t.Date()
	loc := t.Location()
	start = time.Date(y, m-period, 1, 0, 0, 0, 0, loc)
	end = time.Date(y, m, 1, 0, 0, 0, -1, loc)
	return start, end
}

// deleteFile removes a file from the specified Path
func deleteFile(filePath string) (err error) {
	// Check if the file exists before we attempt to delete it
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil
	}

	// delete file
	if err = os.Remove(filePath); err != nil {
		CheckError("Unable to remove File "+filePath, err, false)
		return err
	}

	return nil
}

func renameFile(oldName, newName string) (err error) {
	if err = os.Rename(oldName, newName); err != nil {
		CheckError("Unable to rename File "+oldName, err, false)
		return err
	}

	return nil
}

// SendEmail will send email to given address
func SendEmail(message string, toAddress []string) (response bool, err error) {
	const (
		fromAddress       = "order@slot.ng"
		fromEmailPassword = "=k5qnb45"
		smtpServer        = "smtp.gmail.com"
		smptPort          = "587"
	)
	var auth = smtp.PlainAuth("", fromAddress, fromEmailPassword, smtpServer)
	err = smtp.SendMail(smtpServer+":"+smptPort, auth, fromAddress, toAddress, []byte(message))

	if err == nil {
		return true, nil
	}

	fmt.Println("Mailing error is: ", err)
	return false, err

}

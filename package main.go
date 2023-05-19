package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/badoux/checkmail"
)

type Email struct {
	Email string
	Valid bool
}

type CsvMails struct {
	Emails []*Email
	wg     *sync.WaitGroup
}

func (e *CsvMails) Stats() (int, int) {
	validCount := 0
	for _, email := range e.Emails {
		if email.Valid {
			validCount++
		}
	}
	return validCount, len(e.Emails) - validCount
}

func (e CsvMails) String() string {
	valid, invalid := e.Stats()
	return fmt.Sprintf("Valid: %v\nInvalid: %v\nTotal: %v\n", valid, invalid, len(e.Emails))
}

var (
	serverHostName = "smtp.gmail.com" // set your SMTP server here
	csvFile        = "email1.csv"
	csvFile2       = "email2.csv"
)

func main() {
	var allMails CsvMails
	var wg *sync.WaitGroup
	emails, err := ReadCsvFile(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	emails2, err := ReadCsvFile(csvFile2)
	if err != nil {
		log.Fatal(err)
	}
	allMails.Emails = emails
	allMails.Emails = append(allMails.Emails, emails2...)
	wg.Add(len(allMails.Emails))
	go allMails.ValidateEmail(wg)
}

func ReadCsvFile(filename string) ([]*Email, error) {
	// open file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()

	emails := make([]*Email, 0)

	data, _ := ioutil.ReadFile(filename)
	splits := strings.Split(string(data), "\n")

	for _, line := range splits {
		email := &Email{
			Email: line,
			Valid: false,
		}
		emails = append(emails, email)
	}
	return emails, nil
}

func (e *CsvMails) ValidateEmail(wg *sync.WaitGroup) {
	for _, mail := range e.Emails {
		fmt.Println(mail.Email)
		err := checkmail.ValidateHostAndUser(serverHostName, "", mail.Email)
		if smtpErr, ok := err.(checkmail.SmtpError); ok && err != nil {
			fmt.Printf("Code: %s, Msg: %s", smtpErr.Code(), smtpErr)
		} else {
			mail.Valid = true
		}
	}
	defer e.wg.Done()
}

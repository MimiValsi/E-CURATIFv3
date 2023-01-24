package main

import (
	"bufio"
	_ "encoding/csv"
	"os"
	"os/exec"
	"strings"
)

// We start by verifying if the input file is a .csv
// if true than proceed to csvEncoding()
//
// Execute command to fetch encoding type
//
// Exemple:
//
// file.csv
// file.csv: text/csv; charset=iso-8859-1
//                            ^
//
// 1st run the command with output() to fetch the string
// 2nd split with "="
//
// We'll get:
// str[0] = file.csv: text/csv; charset
// str[1] = iso-8859-1
//
// Copy str[1] to uppercase into a tmp variable
// As we don't know which encoding type a file might have
// for every file scanned we verify it's encoding
// if it's not UTF-8 then run cmd to change to it.

func (app *application) verifyCSV(s string) {
	file := strings.Split(s, ".")
	fileLen := len(file)

	if file[fileLen-1] != "csv" {
		app.infoLog.Print("Wrong type of file")
	} else {
		app.csvEncoding(s)
	}

}

func (app *application) csvEncoding(s string) {

	cmd, err := exec.Command("file", "-i", s).Output()
	if err != nil {
		app.infoLog.Print(err)
	}

	strSplit := []string{}
	tmp := strings.Split(string(cmd), "=")
	strSplit = append(strSplit, tmp...)

	tmp2 := strings.ToUpper(strSplit[1])

	// check if encding type if UTF-8
	// if not than run change encoding cmd
	if tmp2 != "UTF-8" {

		cmd2 := exec.Command("iconv", "-f", tmp2, "-t", "UTF-8",
			s, "-o", s)
		iconvErr := cmd2.Run()
		app.infoLog.Print(iconvErr)
	}


}

func (app *application) checkSeparator(s string) ([]string, error) {
	file, err := os.Open(s)
	if err != nil {
		app.errorLog.Print(err)
	}

	str := []string{}
	fs := bufio.NewScanner(file)
	fs.Split(bufio.ScanLines)

	for fs.Scan() {
		line := fs.Text()
		str = append(str, line)
	}

	return str, nil
}


// TODO:
//
// Verify if in csv, word/characters are separated by ','

/*
   Copyright 2014 Franc[e]sco (lolisamurai@tfwno.gf)
   This file is part of sharenix.
   sharenix is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   sharenix is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with sharenix. If not, see <http://www.gnu.org/licenses/>.
*/

package main

// NOTE: to compile this, you need gtk 2.0 and >=go-1.3.1

import (
	"errors"
	"flag"
	"fmt"

	"github.com/zneix/sharenix/sharenixlib"
)

func handleCLI() (err error) {
	cfg, err := sharenixlib.LoadConfig()
	if err != nil {
		return
	}

	// command line flags
	pmode := flag.String("m", "f",
		"Upload mode - f/file: upload file, fs/fullscreen: screenshot entire "+
			"screen and upload, s/section: select screen region and upload, "+
			"c/clipboard: upload clipboard contents, r/record: record screen "+
			"region and upload, u/url: shorten url")

	psite := flag.String("s", "default",
		"Target site name (default = default site for the selected mode)")

	psilent := flag.Bool("q", false,
		"Quiet mode - disables all terminal output except errors")

	pnotification := flag.Bool("n", false,
		"Notification - displays a GTK notification for the upload")

	popen := flag.Bool("o", false, "Open url - automatically opens the "+
		"uploaded file's url in the default browser")

	pclip := flag.Bool("c", true, "Copy url to clipboard - copies the "+
		"uploaded file's url to the clipboard (not guaranteed to work properly"+
		"on all window managers, tested on Unity + X11)")

	phistory := flag.Bool("history", false, "Show upload history (grep-able)")
	pversion := flag.Bool("v", false, "Shows the program version")
	pdebug := flag.Bool("g", false, "Show verbose debug information "+
		"(output can include sensitive info such as API keys)")

	pupload := flag.Bool("upload", true, "If false, the file will be "+
		"archived but not uploaded")

	flag.Parse()
	if !flag.Parsed() {
		panic(errors.New("Unexpected flag error"))
	}

	if *pversion {
		fmt.Println(sharenixlib.ShareNixVersion)
		return
	}

	if *phistory {
		var csv [][]string
		csv, err = sharenixlib.GetUploadHistory()
		if err != nil {
			return
		}

		if len(csv) == 0 {
			fmt.Println("Empty!")
			return
		}

		for _, record := range csv {
			if len(record) < 4 {
				err = errors.New("Invalid csv")
				return
			}

			fmt.Println("*", record[3], "- URL:", record[0], "Thumbnail URL:",
				record[1], "Deletion URL:", record[2])
			fmt.Println()
		}

		return
	}

	sharenixlib.ShareNixDebug = *pdebug

	// perform upload
	_, _, _, err = sharenixlib.ShareNix(
		cfg, *pmode, *psite, *psilent, *pnotification, *popen, *pclip,
		*pupload)
	if err != nil {
		return
	}

	return
}

func main() {
	err := handleCLI()
	if err != nil {
		fmt.Println(err)
	}
}

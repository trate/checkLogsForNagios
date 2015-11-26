//Parses log during defined last hours and seeks regular expression
//Maintainer: Oleg Laktionov
package main

import (
	"os"
	"fmt"
	"flag"
	"time"
	"log"
	"regexp"
	"io"
	"strings"
	"strconv"
)

func check(e error) {
	if e != nil {
                log.Fatal(e)
        }
}
//For example, from "7" to "07"
func addNull (s string) (snull string) {
	if len(s) == 1 {
		snull = "0"+s
	} else {
		snull = s
	}
	return
}
//Makes slices of hours from past to now
func hoursInterval(hinput, hnow int) (interval []string) {

	var x string
	switch {
	//We want to check a current day only
	case hinput >= hnow: 
		for i := hnow; i >= 0; i-- {
			x = strconv.Itoa(i)
			x = addNull(x)
			interval = append(interval, x)
		}
	//Special case: check only a current hour
	case hinput == 0: 
		i := hnow
                x = strconv.Itoa(i)
                x = addNull(x)
                interval = append(interval, x)
	default:
		hdiff := hnow - hinput
		for i := hnow; i >= hdiff; i-- {
			x = strconv.Itoa(i)
			x = addNull(x)
			interval = append(interval, x)
		}
	}
	return
}
//Make a regular expression string for matching desired time frame
func makeRegexp(interval []string) (sinterval string) {
        var s string
        for _, value := range interval {
                s += "|" + value
        }
        sinterval = s[1:]
        return
}
		

//Parses slice of errors seeking according to provided time values
func readLogsParseTime(b []byte, re, value string, ynow, mnow, dnow, hnow, hinput int) (appErrors []string, ssliceRemain []byte) {

	var sslice []string

        s := string(b)
	

	ynowStr := strconv.Itoa(ynow)
	mnowStr := strconv.Itoa(mnow)
	mnowStr = addNull(mnowStr)
	dnowStr := strconv.Itoa(dnow)
	dnowStr = addNull(dnowStr)

	retimeBreak := regexp.MustCompile(ynowStr +`(/|-)` + mnowStr + `(/|-)` + dnowStr)
// Slicing current day and the rest
	if retimeBreak.MatchString(s) {
		sslice = retimeBreak.Split(s,-1)
		ssliceRemain = []byte(sslice[0])
		sslice = sslice[1:]
	} else {
		nagiosOut(appErrors)	
	}


        retimeint := regexp.MustCompile(`[ T]` + `(` + value + `)` + `:`)
        reerrors := regexp.MustCompile(re)


        l := len(sslice) - 1

// Make slice of the current day only: i > 0, not i >= 0
        for i := l; i >= 0 ; i-- {

                if len(appErrors) == 2 {
                        break
                }

                if retimeint.MatchString(sslice[i]) && reerrors.MatchString(sslice[i]) {
                        appErrors = append(appErrors, sslice[i])
                }
        }

        return
}


func logfile_not_fresh() {
	fmt.Printf("WARNING:  Logfile wasn't modified today!")
	os.Exit(1)
}
//Output for Nagios
func nagiosOut(appErrors []string) {
        if appErrors == nil {
                fmt.Printf("OK: Errors not found.")
                os.Exit(0)
        } else {
		s := appErrors[0]
		s = strings.Replace(s, "\n", " ", -1)
                if len(appErrors) == 1 {
                        fmt.Printf("CRITICAL: %s", s)
                        os.Exit(2)
                } else {
                        fmt.Printf("CRITICAL: Too many errors. Please check logs! %s", s)
                        os.Exit(2)
                }
        }
}


func main() {

	var re string
	var size, i, l, pos int64
	var fEndSlice, ssliceRemain []byte

	flag.StringVar(&re, "r", "[Ee]xception", "A regular expression for seeking")

	hours := flag.Int("h", 1, "Time interval for monitoring in hours\n   Log's date format must be yyyy[/-]mm[/-]dd hh:mm:ss\n   or yyyy[/-]mm[/-]ddThh:mm:ss at the beginning of a line; if --h=0 then it checks a current hour only")
	buf := flag.Int("b", 4096, "Amount of bytes to read from the end of a file")

	flag.Parse()


	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(2)
	}
	
	fname := flag.Arg(0)
	finfo, err := os.Stat(fname)
	if err != nil {
		log.Fatal("Can't open logfile: ", fname)
	}

	mtime := finfo.ModTime()
	tnow := time.Now()

	switch {
	case tnow.Month() != mtime.Month():
		logfile_not_fresh()
	case tnow.Day() != mtime.Day():
		logfile_not_fresh()
	}

	
	hnow := tnow.Hour()
	dnow := tnow.Day()
	mnow := int(tnow.Month())
	ynow := tnow.Year()
	
	

	hinput := *hours
	bs := *buf

	interval := hoursInterval(hinput, hnow)	

        value := makeRegexp(interval)

// File from the end
	f, err := os.Open(fname)
        check(err)

        defer f.Close()

        size, err = f.Seek(0, 2)
        check(err)

	_, err = f.Seek(0, 0)
        check(err)

	//Create buffer equal provided with -b key or by default
	b := make([]byte, bs)
        l = int64(len(b))

        remainder := size % l

	var appErrors []string

	switch {
	// If buffer size is less or equal than file size
	case l <= size:
		for i = l; i <= size; i += l {
			pos = size - i

			_, err = f.Seek(pos, 0)
			check(err)

			_, err = io.ReadFull(f, b)
			check(err)
				
			fEndSlice = b

			//Append a broken part of file (without a date)
			fEndSlice = append(fEndSlice, ssliceRemain...)

			appErrors, ssliceRemain = readLogsParseTime(fEndSlice, re, value, ynow, mnow, dnow, hnow, hinput)

			if len(appErrors) > 0 {
				break
			}
		

		}
		// if something has remained than deal with it
		if remainder != 0 && appErrors == nil {
			b2 := make([]byte, remainder)

			_, err = f.Seek(0, 0)
       			check(err)

			_, err := io.ReadFull(f, b2)
                        check(err)	
			fEndSlice = b2
			fEndSlice = append(fEndSlice, ssliceRemain...)
			appErrors, _ = readLogsParseTime(fEndSlice, re, value, ynow, mnow, dnow, hnow, hinput)
		}
	//If buffer size is greater than a file size
	case l > size:
		_, err = f.Seek(0, 0)
        	check(err)

		b3 := make([]byte, remainder)
                _, err := io.ReadFull(f, b3)
                check(err)
                fEndSlice = b3
                appErrors, _ = readLogsParseTime(fEndSlice, re, value, ynow, mnow, dnow, hnow, hinput)
	}
		

	
	nagiosOut(appErrors)

}

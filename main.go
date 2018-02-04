package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type countryMap struct {
	shortCode string
	country   string
}

var countryMapReal = []countryMap{
	{shortCode: "MUC", country: " DE"},
	{shortCode: "FRA", country: " DE"},
	{shortCode: "ITA", country: "ITA"},
	{shortCode: "HKG", country: "HKG"},
	{shortCode: "AMS", country: " NL"},
	{shortCode: "CH", country: " CH"},
	{shortCode: "CGN", country: " DE"},
	{shortCode: "SZG", country: " AT"},
	{shortCode: "DE", country: " DE"},
	{shortCode: "FR", country: " DE"},
	{shortCode: "GOA", country: "ITA"},
	{shortCode: "SFO", country: "USA"},
	{shortCode: "TFS", country: "ESP"},
	{shortCode: "BER", country: " DE"},
	{shortCode: "US", country: "USA"},
	{shortCode: "Bul", country: "BLG"},
	{shortCode: "FIR", country: "ITA"},
	{shortCode: "AT", country: " AT"},
	{shortCode: "SPA", country: "ESP"},
	{shortCode: "FFM", country: " DE"},
	{shortCode: "ZRH", country: " CH"},
	{shortCode: "VN", country: " VN"},
	{shortCode: "TH", country: " TH"},
}

func countryMapper(shortString string) (countryString string) {

	for _, e := range countryMapReal {
		if e.shortCode == shortString {
			return e.country
		}
	}

	return "error"
}

func parseDate(dateInput string) (timeStamp time.Time) {
	myForm := "20060102T150405Z"
	// t, err := time.Parse(myForm, "20171219T233550Z")
	t, err := time.Parse(myForm, dateInput)
	check(err)
	// if err != nil { fmt.Println(err) }
	fmt.Println("parsed date: ", t)

	// parse Year out of Time
	// y:= t.Year()
	return
}

func readDir() (files []string) {
	root := "./data"
	// root := "./dataSimple"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	// fmt.Println(files)
	// for _, file := range files {
	// 	fmt.Println(file)
	// }
	files = append(files[:0], files[1:]...)
	return
}

func simpleLineCounter() {
	dat, err := ioutil.ReadFile("./sample.txt")
	check(err)
	fmt.Println(string(dat))

	lines := strings.Split(string(dat), "\n")
	fmt.Println(lines)

	for index, element := range lines {
		fmt.Println("line:", index, ", element: ", element)
		// fmt.Printf("%v\n", element)
	}
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)
		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}

func fileIteratorLineCounter(files []string) {
	fmt.Println("\n***** line counter *****")
	for _, f := range files {
		file, _ := os.Open(f)
		// fileScanner := bufio.NewScanner(file)
		r, err := lineCounter(file)
		check(err)
		fmt.Println("lines in file:", r, "filename: ", f)
	}
}

func parseFile(f string) (lines []string, err error) {

	var counter int
	var interimArray []string
	var parserStack [][]string

	dat, err := ioutil.ReadFile(f)
	check(err)
	lines = strings.Split(string(dat), "\n")
	for _, element := range lines {
		// if index < 1 {
		// 	fmt.Println("line:", index, ", element: ", element)
		// }
		interimArray = append(interimArray, element)
		if strings.Contains(element, "BEGIN:VEVENT") {
			counter++
			parserStack = append(parserStack, interimArray)
			interimArray = nil
		}
	}
	fmt.Println("\n****** single file analytics ********")
	fmt.Println("file: ", f)
	fmt.Println("counter BEGIN: ", counter)
	fmt.Println("length of parserStack", len(parserStack))

	// read parserStack

	counterA := 0
	counterB := 0
	counterD := 0
	counterF := 0
	counterG := 0

	for i, e := range parserStack {
		// read one entry
		if i < 200 {
			for _, line := range e {
				if strings.Contains(line, "BEGIN:VEVENT") {
					counterA++
				}
				if strings.Contains(line, "CREATED:") {
					counterB++
				}
				if strings.Contains(line, "DTSTART;VALUE=") {
					counterD++
				}
				if strings.Contains(line, "UID:") {
					counterF++
				}
				if strings.Contains(line, "SUMMARY:") {
					counterG++
				}

			}
		}

	}
	fmt.Println("BEGIN:VEVENT: ", counterA)
	fmt.Println("CREATED: ", counterB)
	fmt.Println("DTSTART;VALUE=:", counterD)
	fmt.Println("UID: ", counterF)
	fmt.Println("SUMMARY: ", counterG)

	return
}

func parseFileIterator(files []string) (data [][]string) {
	for _, f := range files {
		fileData, err := parseFile(f)
		check(err)
		fmt.Println("filename: ", f)
		data = append(data, fileData)
	}
	return
}

type singleElements struct {
	plainEntries []string
	fileName     string
	position     int
	summary      string
	date         string
	year         string
}

type cleanSingleElements struct {
	fileName string
	position int
	summary  string
	country  string
	date     string
	year     string
}

func groupElements(plainList [][]string, files []string) (groupedList []singleElements) {

	fmt.Println("\n******* going into grouping loop *******")

	var plainEntryData []string
	var singleOne singleElements

	for i, e := range plainList {
		for i2, el := range e {
			if strings.Contains(el, "BEGIN:VEVENT") {
				plainEntryData = nil
			}
			if strings.Contains(el, "END:VEVENT") {
				// fmt.Println("end")
				singleOne = singleElements{plainEntries: plainEntryData, fileName: files[i], position: i2, summary: "", date: "", year: ""}
				groupedList = append(groupedList, singleOne)
				plainEntryData = nil
			}
			plainEntryData = append(plainEntryData, el)
		}
	}

	//fmt.Printf("groupedList: %#v\n", groupedList)
	return
}

func cleanAndEnrich(elementList []singleElements) (cleanedGroupedList []singleElements) {
	for _, e := range elementList {
		hit := false
		hitDay := false
		for _, el := range e.plainEntries {
			if strings.Contains(el, "SUMMARY:") {
				if len(el) < 13 {
					// fmt.Println("summary hit:", len(el))
					shortEl := strings.TrimPrefix(el, "SUMMARY:")
					// fmt.Println("summary hit:", shortEl)
					re := regexp.MustCompile(`\r`)
					shortEl = re.ReplaceAllString(shortEl, "")
					e.summary = shortEl
					// fmt.Println("len in:", len(elementList))
					// fmt.Println("len out:", len(cleanedGroupedList))
					hit = true
				}
			}
			if strings.Contains(el, "DTSTART;VALUE=") {
				theDate := strings.TrimPrefix(el, "DTSTART;VALUE=DATE:")
				theDate = theDate[:8]
				theYear := theDate[:4]

				// fmt.Println("date hit:", theDate)
				// fmt.Println("year hit:", theYear)
				e.date = theDate
				e.year = theYear
				hitDay = true
			}
		}
		if hit && hitDay {
			// e.plainEntries = nil
			cleanedGroupedList = append(cleanedGroupedList, e)
			hit = false
			hitDay = false
		}
	}
	fmt.Println("\n******* cleanedGroupedList *******")
	// fmt.Printf("cleanedGroupedList: %#v\n", cleanedGroupedList)
	return
}

func removedPlainData(input []singleElements) (output []cleanSingleElements) {
	var newItem cleanSingleElements
	for _, e := range input {
		// newItem = {}
		newItem.fileName = e.fileName
		newItem.position = e.position
		newItem.summary = e.summary
		newItem.country = ""
		newItem.date = e.date
		newItem.year = e.year
		output = append(output, newItem)
	}
	fmt.Println("\n******* cleanedList *******")
	// fmt.Printf("cleanedSingleList: %#v\n", output)
	return
}

type summaryList struct {
	summary string
	counter int
}

func pullSummaryList(input []cleanSingleElements) {
	var summaryCounter []summaryList
	var hit bool
	for _, e := range input {
		hit = false
		country := e.summary
		for i, el := range summaryCounter {
			if el.summary == country {
				hit = true
				summaryCounter[i].counter++
			}
		}
		if hit == false {
			summaryCounter = append(summaryCounter, summaryList{summary: country, counter: 1})
		}

	}
	fmt.Println("\n******* summaryCounter *******")
	fmt.Printf("summaryCounter: %#v\n", summaryCounter)
}

func updateCountries(input []cleanSingleElements) (output []cleanSingleElements) {
	// countryMapper
	for _, el := range input {
		el.country = countryMapper(el.summary)
		output = append(output, el)
	}
	return
}

type countryList struct {
	country string
	counter int
}

func pullCountryList(input []cleanSingleElements) {
	var countryCounter []countryList
	var hit bool
	for _, e := range input {
		hit = false
		country := e.country
		for i, el := range countryCounter {
			if el.country == country {
				hit = true
				countryCounter[i].counter++
			}
		}
		if hit == false {
			countryCounter = append(countryCounter, countryList{country: country, counter: 1})
		}

	}
	fmt.Println("\n******* countryCounter *******")
	fmt.Printf("countryCounter: %#v\n", countryCounter)
}

type countryYearList struct {
	country string
	year    string
	counter int
}

func pullCountryListByYear(input []cleanSingleElements) {
	var countryCounter []countryYearList
	var hit bool
	for _, e := range input {
		hit = false
		country := e.country
		year := e.year
		for i, el := range countryCounter {
			if el.country == country && el.year == year {
				hit = true
				countryCounter[i].counter++
			}
		}
		if hit == false {
			countryCounter = append(countryCounter, countryYearList{country: country, year: year, counter: 1})
		}

	}
	fmt.Println("\n******* countryCounterByYear *******")
	fmt.Printf("countryCounterByYear: %#v\n", countryCounter)

	sort.Slice(countryCounter, func(i, j int) bool {
		a, _ := strconv.Atoi(countryCounter[i].year)
		b, _ := strconv.Atoi(countryCounter[j].year)
		return a < b
	})

	fmt.Println("\n******* Summary *******")

	for _, e := range countryCounter {
		fmt.Println("year:", e.year, "country:", e.country, "count:", e.counter)
	}
}

func main() {

	files := readDir()

	fileIteratorLineCounter(files)

	plainList := parseFileIterator(files)

	fmt.Println("\n******** plainlist **********")
	fmt.Println("\nlenth of plain array: ", len(plainList))
	for index := 0; index < 7; index++ {
		fmt.Println("", plainList[0][index])
	}

	groupedList := groupElements(plainList, files)

	cleanedGroupedList := cleanAndEnrich(groupedList)

	removedPlainData := removedPlainData(cleanedGroupedList)

	pullSummaryList(removedPlainData)

	countryCleanedData := updateCountries(removedPlainData)

	pullSummaryList(countryCleanedData)

	pullCountryList(countryCleanedData)

	pullCountryListByYear(countryCleanedData)

}

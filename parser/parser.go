package parser

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"strings"
)

type ParsedLineGeneric struct {
	timestamp int
	level     string
	data      string
}

// +gen
type ParsedLine struct {
	Timestamp int
	Level     string
	Hostname  string
	Service   string
	State     string
	CheckType string
	Try       int
	Comment   string
}

func IncludeAny(str string, toInclude []string) bool {
	for _, mustInclude := range toInclude {
		if strings.Index(str, mustInclude) != -1 {
			return true
		}
	}
	return false
}

func shoudParse(line string, mustInclude []string) bool {
	if strings.Contains(line, ":") && IncludeAny(line, mustInclude) {
		return true
	}
	return false
}

func ParseLineGeneric(line string, hosts []string) ParsedLineGeneric {
	if shoudParse(line, hosts) {
		splited := strings.SplitN(line, ":", 2)
		timestamp, err := strconv.Atoi(strings.Trim(strings.Split(splited[0], "]")[0], "["))
		if err != nil {
			return ParsedLineGeneric{}
		}
		level := strings.Trim(strings.Split(splited[0], "]")[1], " ")
		data := strings.Trim(splited[1], " ")
		return ParsedLineGeneric{timestamp, level, data}
	}
	return ParsedLineGeneric{}
}

func ParseLine(line string, hosts []string) ParsedLine {
	preParsed := ParseLineGeneric(line, hosts)
	notParsable := ParsedLineGeneric{}
	if preParsed != notParsable {
		if preParsed.level == "SERVICE ALERT" {
			return ParseLineServiceAlert(preParsed)
		}
		if preParsed.level == "HOST ALERT" {
			return ParseLineHostAlert(preParsed)
		}
		if preParsed.level == "CURRENT HOST STATE" {
			return ParseLineHostAlert(preParsed)
		}
	}
	return ParsedLine{}
}

func ParseLineServiceAlert(preParsed ParsedLineGeneric) ParsedLine {
	splited := strings.Split(preParsed.data, ";")
	if len(splited) != 6 {
		return ParsedLine{}
	}
	try, err := strconv.Atoi(splited[4])
	if err != nil {
		return ParsedLine{}
	}
	hostname := splited[0]
	service := splited[1]
	state := splited[2]
	checkType := splited[3]
	comment := splited[5]
	return ParsedLine{preParsed.timestamp, preParsed.level, hostname, service, state, checkType, try, comment}
}

func ParseLineHostAlert(preParsed ParsedLineGeneric) ParsedLine {
	splited := strings.Split(preParsed.data, ";")
	if len(splited) != 5 {
		return ParsedLine{}
	}
	try, err := strconv.Atoi(splited[3])
	if err != nil {
		return ParsedLine{}
	}
	hostname := splited[0]
	service := ""
	state := splited[1]
	checkType := splited[2]
	comment := splited[4]
	return ParsedLine{preParsed.timestamp, preParsed.level, hostname, service, state, checkType, try, comment}
}

func ParseLog(r io.Reader, hosts []string) ParsedLines {
	var result ParsedLines
	scanner := bufio.NewScanner(r)
	notParsable := ParsedLine{}
	for scanner.Scan() {
		line := scanner.Text()
		parsed := ParseLine(line, hosts)
		if parsed != notParsable {
			result = append(result, parsed)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}

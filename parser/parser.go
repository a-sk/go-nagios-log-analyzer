package parser

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"strings"
)

type ParsedLineGeneric struct {
	timestamp int64
	level     string
	data      string
}

// +gen
type ParsedLine struct {
	Timestamp int64
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
		timestamp, err := strconv.ParseInt(strings.Trim(strings.Split(splited[0], "]")[0], "["), 10, 64)
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
		// Levels that are not covered:
		// * SERVICE NOTIFICATION
		// 	Lines look like
		//	[1452525906] SERVICE NOTIFICATION: slack-channel;srv-dus;Memcached server stats - Pool: name;WARNING;slack-service-notification;MEMCACHE STATS WARNING - hitrate: 30.93 fillrate: 73.24 evictionrate: 370.98 time_get:162 usec time_set:202 usec cur_conn:159;
		//	[1452526252] SERVICE NOTIFICATION: slack-channel;srv-sfo;Mirrormaker lag;WARNING;slack-service-notification;WARNING - [foo:100592 messages lag] [bar:134790 messages lag];
		//	[1452526618] SERVICE NOTIFICATION: daniel;srv-dus;appdata hdfs v6;ACKNOWLEDGEMENT;mail-service-notification;DISK CRITICAL - /appdata/hdfs/v6 is not accessible: Input/output error;Icinga 2 Admin;shi
		// * EXTERNAL COMMAND
		// 	Lines look like
		//	[1452525976] EXTERNAL COMMAND: SCHEDULE_FORCED_SVC_CHECK;srv0;Raid Status;1452525972
		//	[1452526618] EXTERNAL COMMAND: ACKNOWLEDGE_SVC_PROBLEM;srv0;appdata hdfs v6;2;1;0;Icinga 2 Admin;xli
		//	[1452527100] EXTERNAL COMMAND: ENABLE_SVC_CHECK;srv0;transfer_service_Cron
		// * SERVICE FLAPPING ALERT
		// 	Lines look like
		//	[1425428574] SERVICE FLAPPING ALERT: java-cache2-dus;Redis on port 6382;STARTED; Checkable appears to have started flapping (54% change >= 30% threshold)
		//	[1425429669] SERVICE FLAPPING ALERT: imgstore0-dus;Appdata Images Filesystem;STOPPED; Checkable appears to have stopped flapping (28% change < 30% threshold)
		if preParsed.level == "SERVICE ALERT" {
			return ParseLineServiceAlert(preParsed)
		}
		if preParsed.level == "HOST ALERT" {
			return ParseLineHostAlert(preParsed)
		}
		if preParsed.level == "CURRENT HOST STATE" {
			return ParseLineHostAlert(preParsed)
		}
		if preParsed.level == "CURRENT SERVICE STATE" {
			return ParseLineServiceAlert(preParsed)
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

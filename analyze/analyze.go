package analyze

import (
	"github.com/a-sk/go-nagios-log-analyzer/interval"
	"github.com/a-sk/go-nagios-log-analyzer/parser"
	"github.com/deckarep/golang-set"
	"math"
	"regexp"
)

type DownTime struct {
	Hostname          string
	CriticalTimestamp int
	FirstOkTimestamp  int
}

type Uptime map[string]int

func matchAny(str string, regexps []string) bool {
	for _, regexpStr := range regexps {
		r := regexp.MustCompile(regexpStr)
		if r.FindString(str) != "" {
			return true
		}
	}
	return false
}
func CountsAsDown(parsed parser.ParsedLine, criticals []string, tries int) bool {
	if parsed.Try < tries {
		return false
	}
	if parsed.Level == "HOST ALERT" {
		return HostIsDown(parsed)
	}
	if parsed.Level == "SERVICE ALERT" {
		return ServiceIsDown(parsed, criticals)
	}
	return false
}

func HostIsDown(parsed parser.ParsedLine) bool {
	if parsed.State == "DOWN" {
		return true
	}
	return false
}
func ServiceIsDown(parsed parser.ParsedLine, criticals []string) bool {
	if parsed.State != "CRITICAL" {
		return false
	}
	return matchAny(parsed.Service, criticals)
}

func FindDownHosts(parsedLog parser.ParsedLines, criticals []string, tries int) parser.ParsedLines {
	hostsAreDown := func(item parser.ParsedLine) bool {
		return CountsAsDown(item, criticals, tries)
	}
	return parsedLog.Where(hostsAreDown)
}

func index(seq parser.ParsedLines, item parser.ParsedLine) int {
	for idx, el := range seq {
		if item == el {
			return idx
		}
	}
	return -1
}

func FindDowntimeForHost(parsedLog parser.ParsedLines, critical parser.ParsedLine) DownTime {
	criticalIndex := index(parsedLog, critical)
	afterCritical := parsedLog[criticalIndex:]
	firstOk := FindOkCounterPart(afterCritical, critical)
	notFound := parser.ParsedLine{}
	if firstOk != notFound {
		return DownTime{critical.Hostname, critical.Timestamp, firstOk.Timestamp}
	}
	return DownTime{}
}
func FindDowntimeForHosts(parsedLog parser.ParsedLines, criticals []string, tries int) []DownTime {
	downHosts := FindDownHosts(parsedLog, criticals, tries)
	var result []DownTime
	notFound := DownTime{}
	for _, item := range downHosts {
		downtime := FindDowntimeForHost(parsedLog, item)
		if downtime != notFound {
			result = append(result, downtime)
		}
	}
	return result
}

func FindOkCounterPart(parsedLog parser.ParsedLines, critical parser.ParsedLine) parser.ParsedLine {
	if critical.Level == "HOST ALERT" {
		return FindHostIsUp(parsedLog, critical)
	}
	if critical.Level == "SERVICE ALERT" {
		return FindServiceIsOk(parsedLog, critical)
	}
	return parser.ParsedLine{}
}

func FindHostIsUp(parsedLog parser.ParsedLines, critical parser.ParsedLine) parser.ParsedLine {
	for _, currentItem := range parsedLog {
		if currentItem.Hostname == critical.Hostname && currentItem.State == "UP" {
			return currentItem
		}
	}
	return parser.ParsedLine{}

}
func FindServiceIsOk(parsedLog parser.ParsedLines, critical parser.ParsedLine) parser.ParsedLine {
	for _, currentItem := range parsedLog {
		if currentItem.Hostname == critical.Hostname && currentItem.Service == critical.Service && currentItem.State == "OK" {
			return currentItem
		}
	}
	return parser.ParsedLine{}
}

func round(val float64) int {
	return int(math.Floor(val + 0.5))
}

func FindUptime(parsedLog parser.ParsedLines, criticals []string, tries int) Uptime {
	result := make(Uptime)
	tempResult := make(map[string]*interval.IntervalSet)
	for _, downtime := range FindDowntimeForHosts(parsedLog, criticals, tries) {
		_, present := tempResult[downtime.Hostname]
		if !present {
			tempResult[downtime.Hostname] = &interval.IntervalSet{}
		}
		tempResult[downtime.Hostname].Add(interval.Interval{downtime.CriticalTimestamp, downtime.FirstOkTimestamp})
	}
	for hostname, downtimes := range tempResult {
		sumDowntime := downtimes.Len()
		result[hostname] = round(100 - ((float64(sumDowntime) / 86400) * 100))
	}
	for hostname := range AllHostsInLog(parsedLog).Iter() {
		hostname := hostname.(string)
		_, present := result[hostname]
		if !present {
			result[hostname] = 100
		}
	}
	return result
}

func AllHostsInLog(parsedLog parser.ParsedLines) mapset.Set {
	result := mapset.NewSet()
	for _, item := range parsedLog {
		if item.Level == "CURRENT HOST STATE" {
			result.Add(item.Hostname)
		}
	}
	return result
}

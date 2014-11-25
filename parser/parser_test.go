package parser

import (
	//"fmt"
	"reflect"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestParseLineGeneric(t *testing.T) {
	hosts := []string{".test.ru", ".test.com"}
	sampleLine := "[1405209600] CURRENT HOST STATE: host2.test.ru;UP;HARD;1;PING OK - Packet loss = 0%, RTA = 0.52 ms"
	answer := ParsedLineGeneric{1405209600, "CURRENT HOST STATE", "host2.test.ru;UP;HARD;1;PING OK - Packet loss = 0%, RTA = 0.52 ms"}
	parsed := ParseLineGeneric(sampleLine, hosts)

	expect(t, parsed, answer)
}

func TestParseLineGenericNotParsable(t *testing.T) {
	hosts := []string{".test.ru", ".test.com"}
	sampleLine := "[1405209600] LOG VERSION: 2.0"
	answer := ParsedLineGeneric{}
	parsed := ParseLineGeneric(sampleLine, hosts)

	expect(t, parsed, answer)
}

func TestParseLineGenericWrongHost(t *testing.T) {
	hosts := []string{".test.ru", ".test.com"}
	sampleLine := "[1405209600] CURRENT HOST STATE: hello.test3.ru;UP;HARD;1;PING OK - Packet loss = 0%, RTA = 0.52 ms"
	answer := ParsedLineGeneric{}
	parsed := ParseLineGeneric(sampleLine, hosts)

	expect(t, parsed, answer)
}
func TestIncludeAnyTrue(t *testing.T) {
	hosts := []string{".test.ru", ".test.com"}
	sampleLine := "[1405209600] CURRENT HOST STATE: host2.test.ru;UP;HARD;1;PING OK - Packet loss = 0%, RTA = 0.52 ms"
	expect(t, IncludeAny(sampleLine, hosts), true)

}
func TestIncludeAnyFalse(t *testing.T) {
	hosts := []string{"(dedi|vip|vh)[0-9]+\\.test\\.ru", ".\\.test2\\.ru"}
	sampleLine := "[1405209600] CURRENT HOST STATE: somehost.ru;UP;HARD;1;PING OK - Packet loss = 0%, RTA = 0.52 ms"
	expect(t, IncludeAny(sampleLine, hosts), false)

}

func TestParseLineServiceAlert(t *testing.T) {
	sample := ParsedLineGeneric{1405209605, "SERVICE ALERT", "host3.test.ru;SSH;CRITICAL;HARD;3;Server answer:"}
	answer := ParsedLine{1405209605, "SERVICE ALERT", "host3.test.ru", "SSH", "CRITICAL", "HARD", 3, "Server answer:"}
	parsed := ParseLineServiceAlert(sample)

	expect(t, parsed, answer)
}
func TestParseLineServiceAlertGiberish(t *testing.T) {
	sample := ParsedLineGeneric{1405209605, "SERVICE ALERT", "sample string"}
	answer := ParsedLine{}
	parsed := ParseLineServiceAlert(sample)

	expect(t, parsed, answer)
}

func TestParseLine(t *testing.T) {
	hosts := []string{".test.ru", ".test.com"}
	sample := "[1405209605] SERVICE ALERT: host3.test.ru;SSH;CRITICAL;HARD;3;Server answer:"
	answer := ParsedLine{1405209605, "SERVICE ALERT", "host3.test.ru", "SSH", "CRITICAL", "HARD", 3, "Server answer:"}
	parsed := ParseLine(sample, hosts)

	expect(t, parsed, answer)
}
func TestParseLineWithHost(t *testing.T) {
	hosts := []string{".test.ru", ".test.com"}
	sample := "[1405217823] HOST ALERT: host1.test.ru;DOWN;HARD;1;CRITICAL - Plugin timed out after 4 seconds"
	answer := ParsedLine{1405217823, "HOST ALERT", "host1.test.ru", "", "DOWN", "HARD", 1, "CRITICAL - Plugin timed out after 4 seconds"}
	parsed := ParseLine(sample, hosts)

	expect(t, parsed, answer)
}
func TestParseLineWithHostState(t *testing.T) {
	hosts := []string{".test.ru", ".test.com"}
	sample := "[1405209600] CURRENT HOST STATE: host2.test.ru;UP;HARD;1;PING OK - Packet loss = 0%, RTA = 0.47 ms"
	answer := ParsedLine{1405209600, "CURRENT HOST STATE", "host2.test.ru", "", "UP", "HARD", 1, "PING OK - Packet loss = 0%, RTA = 0.47 ms"}
	parsed := ParseLine(sample, hosts)

	expect(t, parsed, answer)
}

func TestParseLineWrongService(t *testing.T) {
	hosts := []string{".test.ru", ".test.com"}
	sample := "[1405209605] ALERT: host3.test.ru;SSH;CRITICAL;HARD;3;Server answer:"
	answer := ParsedLine{}
	parsed := ParseLine(sample, hosts)

	expect(t, parsed, answer)
}

func TestParseLineHostAlert(t *testing.T) {
	sample := ParsedLineGeneric{1405217823, "HOST ALERT", "host1.test.ru;DOWN;HARD;1;CRITICAL - Plugin timed out after 4 seconds"}
	answer := ParsedLine{1405217823, "HOST ALERT", "host1.test.ru", "", "DOWN", "HARD", 1, "CRITICAL - Plugin timed out after 4 seconds"}
	parsed := ParseLineHostAlert(sample)

	expect(t, parsed, answer)
}

func TestParseLineHostAlertGiberish(t *testing.T) {
	sample := ParsedLineGeneric{1405217823, "HOST ALERT", "some random data"}
	answer := ParsedLine{}
	parsed := ParseLineHostAlert(sample)

	expect(t, parsed, answer)
}

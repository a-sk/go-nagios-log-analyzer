package analyze

import (
	//"fmt"
	"../parser"
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

func TestHostIsDown(t *testing.T) {
	parsed := parser.ParsedLine{1405217823, "HOST ALERT", "host2.test.com", "", "DOWN", "HARD", 1, "CRITICAL - Plugin timed out after 4 seconds"}
	expect(t, HostIsDown(parsed), true)
}
func TestHostIsDownIsUp(t *testing.T) {
	parsed := parser.ParsedLine{1405217823, "HOST ALERT", "host2.test.com", "", "UP", "HARD", 1, "CRITICAL - Plugin timed out after 4 seconds"}
	expect(t, HostIsDown(parsed), false)
}

func TestCountsAsDownNotEnoughTries(t *testing.T) {
	parsed := parser.ParsedLine{1405217823, "HOST ALERT", "host2.test.com", "", "DOWN", "HARD", 1, "CRITICAL - Plugin timed out after 4 seconds"}
	criticals := []string{"SSH", "PING"}
	expect(t, CountsAsDown(parsed, criticals, 3), false)
}
func TestCountsAsDown(t *testing.T) {
	parsed := parser.ParsedLine{1405217823, "HOST ALERT", "host2.test.com", "", "DOWN", "HARD", 3, "CRITICAL - Plugin timed out after 4 seconds"}
	criticals := []string{"SSH", "PING"}
	expect(t, CountsAsDown(parsed, criticals, 3), true)
}

func TestServiceIsDown(t *testing.T) {
	parsed := parser.ParsedLine{1405209605, "SERVICE ALERT", "host1.test.com", "SSH", "CRITICAL", "HARD", 3, "Server answer:"}
	criticals := []string{"SSH", "PING"}
	expect(t, ServiceIsDown(parsed, criticals), true)
}

func TestServiceIsWrongService(t *testing.T) {
	parsed := parser.ParsedLine{1405209605, "SERVICE ALERT", "host1.test.com", "SOME", "CRITICAL", "HARD", 3, "Server answer:"}
	criticals := []string{"SSH", "PING"}
	expect(t, ServiceIsDown(parsed, criticals), false)
}

func TestCountsAsDownService(t *testing.T) {
	parsed := parser.ParsedLine{1405209605, "SERVICE ALERT", "host1.test.com", "SSH", "CRITICAL", "HARD", 3, "Server answer:"}
	criticals := []string{"SSH", "PING"}
	expect(t, CountsAsDown(parsed, criticals, 3), true)
}

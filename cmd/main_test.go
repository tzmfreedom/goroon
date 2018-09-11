package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/tzmfreedom/goroon"
)

var isCloudTestCase = []struct {
	input    string
	expected bool
}{
	{"https://example.com", false},
	{"https://example.cybozu.com", true},
	{"https://example.com/cybozu.com", false},
}

func TestIsCloud(t *testing.T) {
	for _, tc := range isCloudTestCase {
		result := isCloud(tc.input)
		if tc.expected != result {
			t.Fatalf("failed. expected = %v, actual = %v", tc.expected, result)
		}
	}
}

var getSessionTestCase = []struct {
	endpoint        string
	cookie          string
	expectedSession string
	expectedError   string
}{
	{"https://example.com", "CBSESSID=123; path=/", "123", ""},
	{"https://example.com", "JSESSIONID=123; path=/", "", "Authentication Failure: not match cookie response"},
	{"https://example.cybozu.com", "JSESSIONID=123; path=/", "123", ""},
	{"https://example.cybozu.com", "CBSESSID=123; path=/", "", "Authentication Failure: not match cookie response"},
	{"https://example.com/cybozu.com", "CBSESSID=123; path=/", "123", ""},
	{"https://example.com/cybozu.com", "JSESSIONID=123; path=/", "", "Authentication Failure: not match cookie response"},
}

func TestGetSession(t *testing.T) {
	for _, tc := range getSessionTestCase {
		actualSession, actualError := getSession(tc.endpoint, &goroon.Returns{Cookie: tc.cookie})
		if tc.expectedSession != actualSession {
			t.Fatalf("failed. expected = %v, actual = %v", tc.expectedSession, actualSession)
		}
		if actualError != nil {
			if tc.expectedError == "" {
				t.Fatalf("failed. expected = %v, actual = %v", tc.expectedError, actualError)
			}
			if actualError.Error() != tc.expectedError {
				t.Fatalf("failed. expected = %v, actual = %v", tc.expectedError, actualError)
			}
		}
	}
}

func TestBeginningOfDay(t *testing.T) {
	now := time.Now()
	day := beginningOfDay(now)
	if now.Year() != day.Year() {
		t.Fatalf("failed. expected = %v, actual = %v", now.Year(), day.Year())
	}
	if now.Month() != day.Month() {
		t.Fatalf("failed. expected = %v, actual = %v", now.Month(), day.Month())
	}
	if now.Day() != day.Day() {
		t.Fatalf("failed. expected = %v, actual = %v", now.Day(), day.Day())
	}
	if day.Hour() != 0 {
		t.Fatalf("failed. expected = 0, actual = %v", day.Hour())
	}
	if day.Minute() != 0 {
		t.Fatalf("failed. expected = 0, actual = %v", day.Minute())
	}
	if day.Second() != 0 {
		t.Fatalf("failed. expected = 0, actual = %v", day.Second())
	}
	if day.Nanosecond() != 0 {
		t.Fatalf("failed. expected = 0, actual = %v", day.Nanosecond())
	}
	if now.Location() != day.Location() {
		t.Fatalf("failed. expected = %v, actual = %v", now.Location(), day.Location())
	}
}

func TestEndOfDay(t *testing.T) {
	now := time.Now()
	day := endOfDay(now)
	if now.Year() != day.Year() {
		t.Fatalf("failed. expected = %v, actual = %v", now.Year(), day.Year())
	}
	if now.Month() != day.Month() {
		t.Fatalf("failed. expected = %v, actual = %v", now.Month(), day.Month())
	}
	if now.Day() != day.Day() {
		t.Fatalf("failed. expected = %v, actual = %v", now.Day(), day.Day())
	}
	if day.Hour() != 23 {
		t.Fatalf("failed. expected = 23, actual = %v", day.Hour())
	}
	if day.Minute() != 59 {
		t.Fatalf("failed. expected = 59, actual = %v", day.Minute())
	}
	if day.Second() != 59 {
		t.Fatalf("failed. expected = 59, actual = %v", day.Second())
	}
	if day.Nanosecond() != 999999 {
		t.Fatalf("failed. expected = 999999, actual = %v", day.Nanosecond())
	}
	if now.Location() != day.Location() {
		t.Fatalf("failed. expected = %v, actual = %v", now.Location(), day.Location())
	}
}

func TestFormatDate(t *testing.T) {
	now := time.Now()
	actual := formatDate(goroon.XmlDate{now})
	expected := fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
	if expected != actual {
		t.Fatalf("failed. expected = %v, actual = %v", expected, actual)
	}

}

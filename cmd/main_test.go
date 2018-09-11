package main

import (
	"testing"

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

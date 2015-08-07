package main

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type Result int

const (
	PASS Result = iota
	FAIL
)

type Report struct {
	Packages []Package
}

type Package struct {
	Name  string
	Time  int
	Tests []*Test
}

type Test struct {
	Name   string
	Time   int
	Result Result
	Output []string
}

var (
	regexStatus = regexp.MustCompile(`^--- (PASS|FAIL): (.+) \((\d+\.\d+)( seconds|s)\)$`)
	regexResult = regexp.MustCompile(`^(ok|FAIL)\s+(.+)\s(\d+\.\d+)s$`)
)

func Parse(r io.Reader) (*Report, error) {
	reader := bufio.NewReader(r)

	report := &Report{make([]Package, 0)}

	// keep track of tests we find
	tests := make([]*Test, 0)
	testMap := make(map[string]*Test)

	// current test
	var test *Test

	// parse lines
	for {
		l, _, err := reader.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		line := string(l)

		if strings.HasPrefix(line, "=== RUN ") {
			// start of a new test
			test = &Test{
				Name:   line[8:],
				Result: FAIL,
				Output: make([]string, 0),
			}

			if _, exist := testMap[test.Name]; !exist {
				testMap[test.Name] = test
				tests = append(tests, test)
			}
		} else if matches := regexResult.FindStringSubmatch(line); len(matches) == 4 {
			// all tests in this package are finished
			report.Packages = append(report.Packages, Package{
				Name:  matches[2],
				Time:  parseTime(matches[3]),
				Tests: tests,
			})

			tests = make([]*Test, 0)
		} else if test != nil {
			if matches := regexStatus.FindStringSubmatch(line); len(matches) == 5 {
				if nTest, ok := testMap[matches[2]]; ok {
					test = nTest
				} else {
					test.Name = matches[2]
				}

				// test status
				if matches[1] == "PASS" {
					test.Result = PASS
				} else {
					test.Result = FAIL
				}

				test.Time = parseTime(matches[3]) * 10
			} else if strings.HasPrefix(line, "\t") {
				// test output
				test.Output = append(test.Output, line[1:])
			}
		}
	}

	return report, nil
}

func parseTime(time string) int {
	t, err := strconv.Atoi(strings.Replace(time, ".", "", -1))
	if err != nil {
		return 0
	}
	return t
}

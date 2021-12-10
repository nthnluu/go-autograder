package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// AutograderConfig is a struct that represents the parsed contents of autograder.config.json
type AutograderConfig struct {
	Visibility string `json:"visibility"`
	Tests      []struct {
		Name       string  `json:"name"`
		Number     string  `json:"number"`
		Points     float64 `json:"points"`
		Visibility string  `json:"visibility,omitempty"`
	} `json:"tests"`
}

// TestResult is a struct that represents the result of a test case in Gradescope's specifications
// https://gradescope-autograders.readthedocs.io/en/latest/specs/
type TestResult struct {
	Score      float64 `json:"score"`
	MaxScore   float64 `json:"max_score"`
	Name       string  `json:"name"`
	Number     string  `json:"number"`
	Output     string  `json:"output"`
	Visibility string  `json:"visibility,omitempty"`
}

// AutograderOutput represents the output that conforms to Gradescope's specifications
// https://gradescope-autograders.readthedocs.io/en/latest/specs/
type AutograderOutput struct {
	Visibility string       `json:"visibility,omitempty"`
	Tests      []TestResult `json:"tests"`
}

// parseTestOutput converts the stdout of running go test into a mapping between test names and whether
// they passed or not.
func parseTestOutput(rawTestOutput []byte) (results map[string]struct {
	Passed bool
	Output string
}) {
	results = make(map[string]struct {
		Passed bool
		Output string
	})

	scanner := bufio.NewScanner(strings.NewReader(string(rawTestOutput)))
	currTest := ""
	currTestOutput := ""

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) >= 7 && line[:7] == "=== RUN" {
			// begin test
			currTest = line[10:]
		} else if len(line) >= 8 && (line[:8] == "--- PASS" || line[:8] == "--- FAIL") {
			// end test and record result
			testPassed := line[4:8] == "PASS"
			results[currTest] = struct {
				Passed bool
				Output string
			}{Passed: testPassed, Output: currTestOutput}

			// reset currTest and currTestOutput
			currTest = ""
			currTestOutput = ""
		} else if currTest != "" {
			currTestOutput += line
		}
	}

	return
}

func JsonTestRunner() (result AutograderOutput, err error) {
	// Open the autograderconfig JSON file
	testConfigPath, err := filepath.Abs("../../autograder.config.json")
	if err != nil {
		return
	}

	file, err := ioutil.ReadFile(testConfigPath)
	if err != nil {
		return
	}

	// Parse the JSON into an array of testConfig structs
	var autograderConfig AutograderConfig
	err = json.Unmarshal(file, &autograderConfig)
	if err != nil {
		return
	}

	// Run all the tests within the submission folder
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	// Change working directory to the student submission
	err = os.Chdir(fmt.Sprintf("%v/../../submission", wd))
	if err != nil {
		return
	}

	// Run go test in the student submission
	out, err := exec.Command("go", "test", "--v", "./...").CombinedOutput()
	if out == nil && err != nil {
		return
	}

	err = nil
	testResults := parseTestOutput(out)

	// Generate autograder output from test results
	result.Visibility = autograderConfig.Visibility
	for _, testConfig := range autograderConfig.Tests {
		testRes, ok := testResults[testConfig.Name]
		if ok {
			res := TestResult{
				Score:      0,
				MaxScore:   testConfig.Points,
				Name:       testConfig.Name,
				Number:     testConfig.Number,
				Visibility: testConfig.Visibility,
			}

			if testRes.Passed {
				res.Score = testConfig.Points
			} else {
				res.Output = testRes.Output
			}

			result.Tests = append(result.Tests, res)
		} else {
			res := TestResult{
				Score:      0,
				MaxScore:   testConfig.Points,
				Name:       testConfig.Name,
				Number:     testConfig.Number,
				Visibility: testConfig.Visibility,
				Output:     "This test failed to run on your submission.",
			}

			result.Tests = append(result.Tests, res)
		}
	}

	return
}

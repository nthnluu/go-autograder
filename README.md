# go-autograder
A boilerplate for building Gradescope autograders for Go projects.

## Getting started
This autograder works by running all Go tests in a student's submission, parsing the results from stdout, and generating a `results.json` file in Gradescope's specified format. Only tests that you configure in `autograder.config.json` are graded and sent to Gradescope. 

When the autograder runs, the student's submission will be copied into `/autograder/source/submission`. You can make any necessary changes to a student's submission -- such as copying in test suite files -- before the autograder runs by adding shell commands to `run_autograder` in the indicated area.

## File hierarchy
- `autograder.config.json` - A JSON file that specifies config options for your autograder suite. Here is where you will define the names of tests and associated point values.
- `setup.sh` - A setup (Bash) script that installs all your dependencies. Gradescope uses Docker running on Ubuntu 18.04 images, so you can use apt, or any other means of setting up packages. By default, it simply uses `apt-get` to install Go.
- `run_autograder` - An executable script, in any language (with appropriate #! line), that compiles and runs your autograder suite and produces the output in the correct place.
- `src/test_runner` - A Go module containing the code responsible for running `go test` on a student's submission, parsing the results from stdout, and returning a `results.json` file in Gradescope's specified format.

## `autograder.config.json`
This JSON file is where you will configure your autograder for your particular assignment. In this file, you must specify the names of the tests you want to use for grading, along with associated point values.

```json=
{
    "visibility": "visible", // Optional visibility setting for autograder results: visible, hidden, after_due_date, after_published
    "tests": [
        {
            "name": "TestAddTwoNumbers",  // The name of the test (must match the test name as defined in test files)
            "number": "1.1", // Optional (will just be numbered in order of array if no number given)
            "points": 5, // The point value of the test case
            "visibility": "visible" // Optional visibility setting for test case: visible, hidden, after_due_date, after_published
        },
        {
            "name": "TestAddTwoNegativeNumbers",
            "number": "1.2",
            "points": 5,
            "visibility": "visible"
        },
        {
            "name": "TestAddNums",
            "number": "2.1",
            "points": 5,
            "visibility": "visible"
        },
        {
            "name": "TestAddNumsOne",
            "number": "2.2",
            "points": 5,
            "visibility": "visible"
        }
    ]
}
```

# go-autograder
A boilerplate for building Gradescope autograders for Go projects.

## Getting started
This autograder works by running all Go tests in a student's submission, parsing the results from stdout, and generating a `results.json` file in Gradescope's specified format. Only tests that you configure in `autograder.config.json` are graded and sent to Gradescope. 

When the autograder runs, the student's submission will be copied into `/autograder/source/submission`. You can make any necessary changes to a student's submission -- such as copying in test suite files -- before the autograder runs by adding shell commands to `run_autograder` in the indicated area.

## File structure overview
- `autograder.config.json` - A JSON file that specifies config options for your autograder suite. Here is where you will define the names of tests and associated point values.
- `setup.sh` - A setup (Bash) script that installs all your dependencies. We're running on Ubuntu 18.04 images, so you can use apt, or any other means of setting up packages. By default, it simply uses `apt-get` to install `go`.
- `run_autograder` - An executable script, in any language (with appropriate #! line), that compiles and runs your autograder suite and produces the output in the correct place.
- `src/test_runner` - A Go module containing the code responsible for running `go test` on a student's submission, parsing the results from stdout, and returning a `results.json` file in Gradescope's specified format.

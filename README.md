# go_exercises

# Exercise #1

# Exercise #2
## Get YAML support package
To get the package, issue `go get gopkg.in/yaml.v3` in a terminal.

## Install the test generator tool
To install the test generator tool, issue `go install github.com/cweill/gotests` in a terminal.

## Create test coverage
First, run test to produce coverage output: `go test exercise_4 -coverprofile=./exercise_4_cover.out`
Second, run tool to create HTML output: `go tool cover -html=./exercise_4_cover.out -o ./exercise_4_cover.html`
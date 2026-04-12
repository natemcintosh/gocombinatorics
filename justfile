default:
    @just --list

alias b := build
alias t := test
alias f := fmt

# Build the project
build:
    go build ./...

# Run tests; optionally specify a test name regex (default: all)
test name=".":
    go test -run "{{name}}" ./...

# Format the code
fmt:
    go fmt ./...

# Run benchmarks; optionally specify a benchmark regex (default: all)
# Regular tests are suppressed via -run='^$'
bench name=".":
    go test -bench="{{name}}" -benchmem -run='^$' ./...

# Run a specific fuzz target (required); e.g. `just fuzz FuzzNewCombinations`
fuzz target:
    go test -fuzz="{{target}}" ./...

# List available fuzz targets
fuzz-list:
    grep -rh --include='*_test.go' '^func Fuzz' . | sed 's/func \(Fuzz[^(]*\).*/\1/'

echo "Testing..."
go test ./... || { echo "Tests failed" ; exit 1; }

echo "Buildling..."
go build

echo "Running Skill Directory..."
./skilldirectory

echo "Testing..."
go test ./controller || { echo "Tests failed" ; exit 1; }
go test ./handler || { echo "Tests failed" ; exit 1; }
go test ./model || { echo "Tests failed" ; exit 1; }
go test ./data || { echo "Tests failed" ; exit 1; }

echo "Buildling..."
go build

echo "Running Skill Directory..."
./skilldirectory

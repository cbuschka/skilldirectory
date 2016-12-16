echo "Testing..."
go test ./controller || { echo "Tests failed" ; exit 1; }
go test ./handler || { echo "Tests failed" ; exit 1; }
go test ./model || { echo "Tests failed" ; exit 1; }
go test ./data || { echo "Tests failed" ; exit 1; }

echo "Buildling..."
go build

running=$(docker inspect -f {{.State.Running}} cassandra_container)
echo $running
if $running; then
  echo "Already Running"
else
  docker-compose build
  docker-compose up -d
  sleep 20
  docker exec -it cassandra_container bash usr/bin/cqlsh -f /data/skilldirectoryschema.cql

fi

echo "Running Skill Directory..."
./skilldirectory

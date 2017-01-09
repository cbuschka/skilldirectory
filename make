################################################################################
# This file is used to run SkillDirectory's unit tests, build the project's    #
# executable, setup Docker containers to hold the Cassandra database and to    #
# run SkillDirectory itself, and lastly, run the executable in the Ubuntu      #
# container.                                                                   #
################################################################################

### Default flags and env vars
export CASSANDRA_USERNAME=cassandra
export CASSANDRA_PASSWORD=cassandra
drop_data_flag=false
export DEBUG_FLAG=true

### Parse all command line flags
for arg in "$@"
do
  if [[ $arg = "--dropdata" ]]; then
    drop_data_flag=true
  elif [[ $arg = "--nodebug" ]]; then
    export DEBUG_FLAG=false 
  else
    echo Unrecognized option: \"$arg\"
    echo Valid options are: \"--dropdata\" and \"--nodebug\"
    exit 127 # exit code for option not found
  fi
done

### Run project tests with 'go test'
echo "Running Tests..."
go test $(glide novendor) || { echo "Tests failed" ; exit 1; }

### Build executable for Ubuntu docker container
env GOOS=linux GOARCH=amd64 go build

### Build the Cassandra and SkillDirectory Docker images
docker-compose build

### See if containers for Cassandra and/or Skilldirectory are cassandra_running
cassandra_running=$(docker inspect -f {{.State.Running}} cassandra_container)
skilldirectory_running=$(docker inspect -f {{.State.Running}} skilldirectory_container)

### Restart/start the skilldirectory_container
if $skilldirectory_running; then
  echo 'Stopping skilldirectory_container...'
  docker stop skilldirectory_container >/dev/null 
  echo 'skilldirectory_container stopped.'
  skilldirectory_running=false
else
  echo 'skilldirectory_container already stopped.'
fi

### If cassandra container is running and "--dropdata" flag was used, stop the container
if $cassandra_running && $drop_data_flag; then
    echo 'Stopping cassandra_container...'
    docker stop cassandra_container >/dev/null
    echo 'cassandra_container stopped.'
    cassandra_running=false
fi

### Get the cassandra_container built and running if it isn't already
if $cassandra_running; then
  echo "cassandra_container is already cassandra_running"
else
  docker-compose up -d cassandra
  sleep 20
fi

### If "--dropdata" flag was used, drop the project's Cassandra keyspace within container
if $drop_data_flag; then
    echo 'Dropping and rebuilding "skill_directory_keyspace" keyspace'
    docker exec -it cassandra_container bash usr/bin/cqlsh -u $CASSANDRA_USERNAME -p $CASSANDRA_PASSWORD -e "DROP KEYSPACE skill_directory_keyspace"
fi

### Execute CQL commands in the container from schema file to set up database
echo "Running skilldirectoryschema..."
docker exec -it cassandra_container bash usr/bin/cqlsh -u $CASSANDRA_USERNAME -p $CASSANDRA_PASSWORD -f /data/skilldirectoryschema.cql
echo "Schema update complete."

### Start Ubuntu container - will run skilldirectory executable on startup and
### connect to cassandra_container for database connectivity
echo "Running Skill Directory..."
docker-compose up skilldirectory # Run with debug-level logging, unless "--nodebug" flag was used

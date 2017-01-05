################################################################################
# This file is used to run SkillDirectory's unit tests, build the project's    #
# executable, setup a Docker container holding the Cassandra database used by  #
# SkillDirectory, and lastly, run the executable                               #
################################################################################

### Parse all command line flags
drop_data_flag=false
debug_flag=true
for arg in "$@"
do
  if [[ $arg = "--dropdata" ]]; then
    drop_data_flag=true
  elif [[ $arg = "--nodebug" ]]; then
    debug_flag=false 
  else
    echo Unrecognized option: \"$arg\"
    echo Valid options are: \"--dropdata\" and \"--nodebug\"
    exit 127 # exit code for option not found
  fi
done

### Run project tests with 'go test'
echo "Running Tests..."
go test $(glide novendor) || { echo "Tests failed" ; exit 1; }

### Create environment vars used by project
export CASSANDRA_URL='0.0.0.0'
export CASSANDRA_PORT=''
export CASSANDRA_KEYSPACE='skill_directory_keyspace'

### Build the project's executable
go build

### See if the docker container is running
running=$(docker inspect -f {{.State.Running}} cassandra_container)

### If container is running and "--dropdata" flag was used, stop the container
if $running && $drop_data_flag; then
    echo 'Stopping cassandra_container...'
    docker stop cassandra_container >/dev/null
    echo 'cassandra_container stopped.'
    running=false
fi

### Get the container built and running if it isn't already
if $running; then
  echo "cassandra_container is already running"
else
  docker-compose build
  docker-compose up -d
  sleep 20
fi

### If "--dropdata" flag was used, drop the project's Cassandra keyspace within container
if $drop_data_flag; then
    echo 'Dropping and rebuilding "skill_directory_keyspace" keyspace'
    docker exec -it cassandra_container bash usr/bin/cqlsh -e "DROP KEYSPACE skill_directory_keyspace"
fi

### Execute CQL commands in the container from schema file to set up database
echo "Running skilldirectoryschema..."
docker exec -it cassandra_container bash usr/bin/cqlsh -f /data/skilldirectoryschema.cql
echo "Schema update complete"

### Run SkillDirectory executable - note that the executable isn't run in the docker
### container. The container only holds the Cassandra database that the executable
### connects to.
echo "Running Skill Directory..."
./skilldirectory -debug=$debug_flag # Run with debug-level logging, unless "--nodebug" flag was used

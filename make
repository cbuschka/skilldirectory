# See if the "--dropdata" flag was used
drop_data_flag=false
if [[ ! -z $1 ]]; then
    if [[ $1 = "--dropdata" ]]; then
        drop_data_flag=true
    elif [[ $1 = "--debug" ]]; then
        drop_data_flag=true
    else
         echo Unrecognized option: \"$1\"
         echo Did you mean \"--dropdata?\"
         exit 127
    fi
fi

echo "Testing..."
go test $(glide novendor)|| { echo "Tests failed" ; exit 1; }

export CASSANDRA_URL='0.0.0.0'
export CASSANDRA_PORT=''
export CASSANDRA_KEYSPACE='skill_directory_keyspace'
export CASSANDRA_USERNAME='cassandra'
export CASSANDRA_PASSWORD='cassandra'

# echo $CASSANDRA_URL
# echo $CASSANDRA_KEYSPACE
# echo "Buildling..."
go build

# docker-compose up -d

running=$(docker inspect -f {{.State.Running}} cassandra_container)
if $running && $drop_data_flag; then
    echo 'Stopping cassandra_container'
    docker stop cassandra_container >/dev/null
    running=false
fi

if $running; then
  echo "cassandra_container is already running"
else
  docker-compose build
  docker-compose up -d
  sleep 20

fi

# If "--dropdata" flag was used, then drop the keyspace
if $drop_data_flag; then
    echo 'Dropping and rebuilding "skill_directory_keyspace" keyspace'
    docker exec -it cassandra_container bash usr/bin/cqlsh -e "DROP KEYSPACE skill_directory_keyspace"
fi

echo "Running skilldirectoryschema"
docker exec -it cassandra_container bash usr/bin/cqlsh -f /data/skilldirectoryschema.cql
echo "Schema update complete"

echo "Running Skill Directory..."
./skilldirectory

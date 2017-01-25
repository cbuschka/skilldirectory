################################################################################
# This file is used to run SkillDirectory's unit tests, build the project's    #
# executable, setup Docker containers to hold the Cassandra database and to    #
# run SkillDirectory itself, and lastly, run the executable in the Ubuntu      #
# container.                                                                   #
################################################################################

### Default flags and env vars
export CASSANDRA_USERNAME=cassandra
export CASSANDRA_PASSWORD=cassandra
export CASSANDRA_URL="0.0.0.0"
export CASSANDRA_PORT=9042
export CASSANDRA_KEYSPACE=skill_directory_keyspace
export DEBUG_FLAG=true
export FILE_SYSTEM=LOCAL

### Parse all command line flags
for arg in "$@"
do
  if [[ $arg = "--nodebug" ]]; then
    export DEBUG_FLAG=false
  else
    echo Unrecognized option: \"$arg\"
    echo Valid options are: \"--nodebug\"
    exit 127 # exit code for option not found
  fi
done

### Run project tests with 'go test'
#echo "Running Tests..."
#go test $(glide novendor) || { echo "Tests failed" ; exit 1; }

echo "Running skilldirectory project..."
### Build executable for Ubuntu docker container
go run main.go -debug=$DEBUG_FLAG

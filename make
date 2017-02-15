################################################################################
# This file is used to run SkillDirectory's unit tests, build the project's    #
# executable, and run it.                                                      #
################################################################################

### Default flags and env vars
export POSTGRES_USERNAME=postgres
export POSTGRES_PASSWORD=password
export POSTGRES_URL="0.0.0.0"
export POSTGRES_PORT=5432
export POSTGRES_KEYSPACE=skilldirectory
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
echo "Running Tests..."
go test $(glide novendor) || { echo "Tests failed" ; exit 1; }

echo "Running skilldirectory project..."
go run main.go -debug=$DEBUG_FLAG

[![Stories in Ready](https://badge.waffle.io/maryvilledev/skilldirectory.png?label=ready&title=Ready)](http://waffle.io/maryvilledev/skilldirectory)
[![CircleCI](https://circleci.com/gh/maryvilledev/skilldirectory.svg?style=svg)](https://circleci.com/gh/maryvilledev/skilldirectory)
[![Report](https://goreportcard.com/badge/github.com/maryvilledev/skilldirectory)](https://goreportcard.com/report/github.com/maryvilledev/skilldirectory)


# SkillDirectory
SkillDirectory is a REST API written in Go.
Our current line of thought regarding this project is that it will be used to
store information about our dev team: teammembers, skills they have, web links
for improving one's knowledge of a subject, etc... Since SkillDirectory is a 
REST API, from the point of view of a client, all of this information is 
created, updated, deleted, and otherwise modified via HTTP requests to a server 
running an instance of SkillDirectory. The server itself stores this data in a
[Cassandra](http://cassandra.apache.org/) database (though this is irrelevant to
clients).  

We are working on [a frontend client using ReactJS](https://github.com/maryvilledev/skilldirectoryui).

## Running SkillDirectory
**1)** Make sure Golang is installed, and your GOPATH is set 
(see https://golang.org/doc/install).

**2)** Install glide by running `brew install glide` (You'll need to have 
[Homebrew](http://brew.sh/) installed to do that).

**3)** Install `http-server` by running `npm install http-server -g`
(you will need [`npm`](https://nodejs.org/en/) installed to do that).

**4)** Clone the repo onto your machine:
```
cd $GOPATH/src # MUST clone into this specific directory!
git clone https://github.com/maryvilledev/skilldirectory.git
```

**5)** Install/download the dependencies:
```
cd $GOPATH/src/skilldirectory
glide install
```

**6)** Run `./make` to start the API server locally on port `8080`
(see https://github.com/maryvilledev/skilldirectory/wiki/Make-File for more info 
on the `make` file).

**Note: This repository only contains the REST API layer, and does not contain the
database that it depends upon. Please make sure you have the Cassandra database from 
[skilldirectoryinfra](https://github.com/maryvilledev/skilldirectoryinfra) running, 
and listening on port `9042`, or the API will not be able to do anything.**

Please also read the [REST Requests](https://github.com/maryvilledev/skilldirectory/wiki/REST-Requests) 
wiki page to learn how to interact with SkillDirectory once you've got it running.

## Wiki
Please explore the project [Wiki](https://github.com/maryvilledev/skilldirectory/wiki)
for more information.

## Endpoints
Below is a listing of all endpoints supported by the API:

* `/links`
* `/skillreviews`
* `/skills`
* `/teammembers`
* `/tmskills`
* `/skillicons`

[![Stories in Ready](https://badge.waffle.io/maryvilledev/skilldirectory.png?label=ready&title=Ready)](http://waffle.io/maryvilledev/skilldirectory)
[![CircleCI](https://circleci.com/gh/maryvilledev/skilldirectory.svg?style=svg)](https://circleci.com/gh/maryvilledev/skilldirectory)
[![Report](https://goreportcard.com/badge/github.com/maryvilledev/skilldirectory)](https://goreportcard.com/report/github.com/maryvilledev/skilldirectory)


# SkillDirectory
Welcome to the SkillDirectory project! This is what you will be working on for
at least your first few weeks at Maryville. The project is relatively new (The
initial commit was made on Nov 28th, 2016), and is very much in flux. Currently 
this is where we are at:

SkillDirectory is a REST API written in Go. If you're unfamiliar with REST, it
is essentially a method by which a client and server can communicate through
HTTP requests and responses. We are also working on [a frontend using ReactJS](https://github.com/maryvilledev/skilldirectoryui).

Our current line of thought regarding this project is that it will be used to
store information about our dev team: teammembers, skills they have, web links
for improving one's knowledge of a subject, etc... Since SkillDirectory is a 
REST API, from the point of view of a client, all of this information is 
created, updated, deleted, and otherwise modified via HTTP requests to a server 
running an instance of SkillDirectory. The server itself stores this data in a
[Cassandra](http://cassandra.apache.org/) database (though this is irrelevant to
clients).  

You'll notice that there are a few badges at the top of this README:

* [![Stories in Ready](https://badge.waffle.io/maryvilledev/skilldirectory.png?label=ready&title=Ready)](http://waffle.io/maryvilledev/skilldirectory)
This badge will take you to the Waffle board for the SkillDirectory project.
We use Waffle to track the project's issues/stories, and track who is working on
what, what's been completed, etc... You can do everything in GitHub that you can
do in Waffle, but you don't get the nice visual representation of the project 
that way.

* [![CircleCI](https://circleci.com/gh/maryvilledev/skilldirectory.svg?style=svg)](https://circleci.com/gh/maryvilledev/skilldirectory)
This badge will take you to the CircleCI page for the SkillDirectory project.
We use CircleCI to run our tests and make sure the project builds correctly
whenever a new commit is made, or a branch is merged.


* [![Report](https://goreportcard.com/badge/github.com/maryvilledev/skilldirectory)](https://goreportcard.com/report/github.com/maryvilledev/skilldirectory)
This badge will take you to the Go Report Card page for the SkillDirectory
project. We use the Go Report Card to make sure we are writing good Go code and
adhering to the proper standards.

Please explore the project [Wiki](https://github.com/maryvilledev/skilldirectory/wiki) for more information. Please also make
sure to read the [New Hire List](https://github.com/maryvilledev/skilldirectory/wiki/New-Hire-List) page.

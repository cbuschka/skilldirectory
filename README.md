[![Stories in Ready](https://badge.waffle.io/maryvilledev/skilldirectory.png?label=ready&title=Ready)](http://waffle.io/maryvilledev/skilldirectory)
[![CircleCI](https://circleci.com/gh/maryvilledev/skilldirectory.svg?style=svg)](https://circleci.com/gh/maryvilledev/skilldirectory)


##### To run locally:

1. Install Golang and set your GOPATH
2. Install glide
  - brew install glide
3. Clone the repo into a `skilldirectory` directory in your go path `go/src/skilldicectory`
4. Add a `skills` folder to the root directory
5. Run `glide install`
6. Run `./make` to run unit test and start server


##### To add a new dependency

1. Add the url to the the list of imports in a file
2. Run `glide config-wizard` to pull dependencies into the `glide.yaml` and `glide.lock`
3. Run `glide install` to add the code to the `vender` directory

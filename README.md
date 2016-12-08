[![Stories in Ready](https://badge.waffle.io/maryvilledev/skilldirectory.png?label=ready&title=Ready)](http://waffle.io/maryvilledev/skilldirectory)
[![CircleCI](https://circleci.com/gh/maryvilledev/skilldirectory.svg?style=svg)](https://circleci.com/gh/maryvilledev/skilldirectory)


### To run locally:

1. Install Golang and set your GOPATH
2. Install glide
  - brew install glide
3. Clone the repo into a `skilldirectory` directory in your go path `go/src/skilldicectory`
4. Add a `skills` folder to the root directory
5. Run `glide install`
6. Run `./make` to run unit test and start server


### To add a new dependency

1. Add the url to the the list of imports in a file
2. Run `glide config-wizard` to pull dependencies into the `glide.yaml` and `glide.lock`
3. Run `glide install` to add the code to the `vender` directory

### To Perform a POST request:
We can add new skills to the SkillDirectory project by sending POST requests to the server running locally.
  1. **Download** [Postman](https://www.getpostman.com/)
  2. **Open** a new tab/request
  3. **Select** "POST" from the dropdown
  4. **Enter** the following URL: `http://localhost:8080/skills`
  5. **Click** the "Body" tab
  6. **Select** the "raw" option
  7. **Select** the "JSON (application/json)" option from the orange dropdown (default is "Text")
  8. **Enter** the following into the text field: `{"Name":"Java","SkillType":"database"}`
  9. **Click** the blue "Send" button
After sending the POST request, you will see console output from the terminal running the SkillDirectory project. In this case, the output would simply be:
```
2016/12/07 15:15:19 Handling Skills Request: POST
2016/12/07 15:15:19 New skill saved
```
Of course, the date and time will vary.

### To Perform a GET Request
We can also perform queries, to obtain a list of skills that have been entered into the SkillHandler database:
  1. **Open** a new Postman tab
  3. **Select** "GET" from the dropdown
  4. **Enter** the following URL: http://localhost:8080/skills/
  9. **Click** the blue "Send" button
After sending the GET request, you will see the following console output:
```
2016/12/07 15:18:56 Handling Skills Request: GET
```
You will also see the following content in the "Body" tab of your GET Request in Postman:
```
[{"Id":"bbc23f7e-bcc2-11e6-9f43-6c4008bcfa84","Name":"Java","SkillType":"database"}]
```
Note that the specific ID may vary. Also note that this is the result of performing a GET request on a SkillDirectory project that has only received a single
POST request. If multiple POST requests have been made, you will see another entry for each additional skill. For example, if you had performed 3 POST requests (and thus
entered 3 skills total), you might see:
```
[{"Id":"59317629-bcc3-11e6-9f43-6c4008bcfa84","Name":"C++","SkillType":"database"},
 {"Id":"62b79945-bcc3-11e6-9f43-6c4008bcfa84","Name":"Golang","SkillType":"database"},
 {"Id":"bbc23f7e-bcc2-11e6-9f43-6c4008bcfa84","Name":"Java","SkillType":"database"}]
```
If you perform a GET request on a SkillDirectory project that has not received any POST requests, then you will only see a pair of square brackets:
```
[]
```

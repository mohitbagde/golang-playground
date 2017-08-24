## CircleCI Build Status : [![CircleCI](https://circleci.com/gh/mohitbagde/golang-playground.svg?style=svg)](https://circleci.com/gh/mohitbagde/golang-playground)


# golang-playground
My go playground to test stuff out on!

### v1.0.0 (08/23/17)
- Added OAuth validation.
- Added basic UI, form templating for the OAuth form page.
- Added a server and handler interface for testing OAuth
- Updated makefile with new commands:
```
make version : Displays the current git commit (or) a "dirty" version if changes were made to the recent branch
make test : Runs the tests for each file
make run : Starts up the binary executable on a server instance that can be cURL'd to
make install : Builds a binary executable for the web app to run
make vet : Vetting the code
make format : Formats the code
make lint : Checks for lint errors (gometalinter is used)
make init : Initializes by fetching any external dependecies/packages
```
todo: 
- Complete coverage (with tests)
- Add a circle yaml file for circleCI to run for every PR and finish up unit tests
- (optional) allow for tokenSecrets to be included in OAuth validation
---


# Madelyne

Tired to use postman to test your API ? 
Have you ever written a curl command to see if your app is responding correctly ?  

Madelyne will help you with that. 
She will run curl-like test on your REST API. 

## Table of contents

- [Madelyne](#madelyne)
  - [Table of contents](#table-of-contents)
  - [Get Madelyne](#get-madelyne)
  - [Usage](#usage)
  - [Config file](#config-file)
  - [Test files](#test-files)
  - [Advanced options](#advanced-options)
  - [Example project](#example-project)
  - [Running tests](#running-tests)
  - [Road map](#road-map)

## Get Madelyne

You can download the latest build directly on [github](https://github.com/madelyne-io/madelyne/releases)

Or you can build the latest version from source

```bash
go get github.com/madelyne-io/madelyne
```

## Usage:

To run the tool, you should just provide the main config file.
The tests are run in the current folder. Be carefull where you are.

```bash
madelyne conf.yml
```

## Config file
The purpose of the config file is to explain to Madelyne what she must do.

```yml
# conf.yml
url: http://localhost:3000
groups:
  main:
    globalSetupCommand: ./example& sleep 1;
    globalTearDownCommand: pkill example
    setupCommand: curl http://localhost:3000/_reset
    teardownCommand: ~
    environment: env.json
    tests: 
      - tests.yml
```

You can cut your tests in group. Like group of feature or by behavior (all the 200 in the same folder), or whatever you want.
Each group could have diferent :

 * Setup and teardown commands which will allow you to load diferent set of fixtures
 * Environment file for more flexibility. See advanced usage for that
 * Set of test files describing all your unit tests and your scenarios. See the next part for that

## Test files

Must be put in a `{groupname}/configs` folder to be found by Madelyne.

To describe your unit test and scenarios. Here what the test file look like

```yaml
# main/configs/tests.yml
unit_tests:
    GET:
        - { url: "/items"                 , status: 200, out: "response/all" }
        - { url: "/items/1/attachments/1" , status: 200, out: "response/file.pdf", ct_out: "application/pdf" }
        - { ... }
    POST:
        - { url: "/item"                  , status: 201, out: 'response/posted', in: 'payload/topost'}
        - { url: "/items/1/attachment"    , status: 201, out: 'response/posted', in: 'payload/file.pdf', ct_in: "application/pdf" }
scenarios:
    scenario1:
        - { action: "POST",  url: "/item",    stauts: 201, in: 'payload/topost' }
        - { action: "GET",   url: "/items/1", stauts: 200, out: "response/one" }
        - { action: "DELET", url: "/items/1", stauts: 204 }
        - { action: "GET",   url: "/items/1", stauts: 404 }
        - { ... }
    scenario2:
        - { ... }
```
The difference between unit test and scenario is when the setup and teardown commands are called.

### Unit tests process
```
globalSetupCommand
setupCommand
unittest1
teardownCommand
setupCommand
unittest2
teardownCommand
globalTearDownCommand
```

### Scenario process

```
globalSetupCommand
setupCommand
scenario1_unittest1
scenario1_unittest2
teardownCommand
setupCommand
scenario2_unittest1
scenario2_unittest2
teardownCommand
globalTearDownCommand
```
The accepted HTTP methods are GET, POST ,PATCH, PUT and DELETE.

Here are the parameters you can provide for any unit test:
|Parameter|Purpose|
|--|--|
|`url`| The url to call. **mandatory** : a test has no meanning  without it. |
|`status`| The expected returned status (default value is 200)|
|`headers`| Headers to send (separated by a `;`) formated by the classic `name : value`|
|`ct_in`| Content-type of what you send (default value is `application/json`)|
|`in`| relative  path to the Content you send from `{groupname}/payloads` folder. If `ct_in` is `application/json` the extension `.json` is added to your filename|
|`ct_out`| Expected content-type  (default value is `application/json`)|
|`out`|  relative  path to the Expected response Content from `{groupname}/responses` folder . If `ct_out` is `application/json` the extension `.json` is added to your filename |

In scenarios, parameters are the same, you just need to provide a `action` parameter (GET, POST, ...)

## Advanced options

If your response can vary and you whant to validate the structure more than the data, then you should read the [advanded option documentation](advanced_readme.md)

## Example project

In the folder `_example` you will find some tests for a fake API as an example. You can read it to get started.

To run it you will required go to be installed

```bash
git clone https://github.com/madelyne-io/madelyne
cd madelyne/_example
go build
mv example tests/example
cd tests 
madelyne conf.yml
```

This should print something like that 

```bash
Testing REST API with Madelyne
[..................................................]100%        6/6
Success
```

## Running tests

```
go test ./... -coverprofile=coverage.out && ./e2e-test.sh
```

## Road map

Next steps are : 

 - reload env.json before each test ?? 
 - add a new command to create default test folder
 - A better output when a test failed, to allow you to localise the test responsible more quickly
 - Add a build fixture phase


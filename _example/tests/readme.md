# Testing example App

## 1 Build app 

You need to build the tester and the app you want to test.
Starting in this folder enter those command

```bash
cd ..
go build
mv example tests/example
cd .. 
go build 
mv madelyne _example/tests/tester
```

## 2 Running the tests

You must be in the tests folder 

```bash 
./madelyne conf.yml
```
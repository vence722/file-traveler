# file-traveler
A small tool for sending a single file to another machine.

## Build
```
go build
```

## Usage
For receiving file (this will start the server listening port 2125):
```
./file-traveler
``` 

For sending a file to a target machine:
```
./file-traveler <file-path> <target-host-name>
```
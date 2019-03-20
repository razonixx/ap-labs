Lab 3.1 - File/Directory Monitor
========================

Description
--------------------------
This is a file and directory monitor written in C using the `inotify` API. This monitor tracks file creation, deletion and renaming.

Compilation
--------------------
To build the binary, use the included `Makefile` by running `make`.

The program recieves 2 arguments:
-  `stdout` or `syslog`, to denote the place where messages will be 
displayed. 

- The folder which will be monitored.


Usage: 
- `./monitor.out stdout .`
- `./monitor.out syslog ./<dir>`


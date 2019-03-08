Lab 2.4 - Syslog Logger Library
========================

Description
--------------------------
This is a library that will print output to the terminal or to the system log depending on which function is used. The library supports: 

- General information with `infof` function.
- Warning messages with `warnf` function.
- Error messages with `errorf` function.
- Error messages and core dump with `panicf` function.
- To choose between `STDOUT` and `SYSLOG`, the function `initLogger` must be used. 
- Calling `initLogger("syslog")` will initialize the log to output the logs to `SYSLOG`. 
- Calling `initLogger("stdout")` or `initLogger("")` will initialize the log to output the logs to `STDOUT`. 
- If `initLogger` is not used, the logger will default to `STDOUT`.

Compilation
--------------------
To use the library, include the `logger.h` header file in your code.

When compiling: include the `logger.c` file in your gcc command. 

To build the example, I've included a Makefile. You can just run `make` to build and `./testLogger.out` to run the test program.
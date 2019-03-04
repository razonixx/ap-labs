Lab 2.2 - Logger Library
========================

Description
--------------------------
This is a library that will print output to the terminal depending on which function is used. The lirary supports: 

- General information with `infof` function.
- Warning messages with `warnf` function.
- Error messages with `errorf` function.
- Error messages and core dump with `panicf` function.

Compilation
--------------------
To use the library, include the `logger.h` header file in your code.

When compiling: include the `logger.c` file in your gcc command. 

To build the example, I've included a Makefile. You can just run `make` to build and `./testLogger.out` to run the test program.
Lab 4.2 - Matrix Multiplication
===============================

Description
--------------------------------
On this lab, I created a multi-threaded matrix multiplicator with the usage of the `pthreads` library. I applied the following concepts:

- Multitheading
- Synchronization mechanisms

The program needs 2 files, each representing one matrix. Each file must contain one number per line. Currently, only `2000x2000` matrices are supported.

The program emulates limited memory, by restricting the amount of data that can be manipulated at any given time with the use of `mutex`.

Compilation
---------------------------------
To build the binary, use the included `Makefile` by running `make`.

The program recieves one parameter, the number of buffers you would like to use. The amount of buffers represents the number of concurrent dot product operations that can be executed at any given time.

Usage
-----------
- `./multiplier -n <buffers>`

Example
- `./multupler -n 10`

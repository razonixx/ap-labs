CC=gcc

#By putting the object files in the rule, make knows it must first compile the .c versions individually and then build the executable.

strlib: strlib.o str-main.o
	$(CC) -o str-final strlib.o str-main.o
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <execinfo.h>

#define RESET		0
#define BRIGHT 		1
#define DIM			2
#define UNDERLINE 	3
#define BLINK		4
#define REVERSE		7
#define HIDDEN		8

#define BLACK 		0
#define RED			1
#define GREEN		2
#define YELLOW		3
#define BLUE		4
#define MAGENTA		5
#define CYAN		6
#define	WHITE		7

void textcolor(int attr, int fg, int bg)
{	char command[13];

	/* Command is the control command to the terminal */
	sprintf(command, "%c[%d;%d;%dm", 0x1B, attr, fg + 30, bg + 40);
	printf("%s", command);
}	

void print_backtrace(void){

	void *tracePtrs[20];
	size_t count;
	count = backtrace(tracePtrs, 10);
	char** funcNames = backtrace_symbols(tracePtrs, count);
	printf("\nBEGIN CORE DUMP\n");
	printf("-----------------------------------------------------\n");
	for (int i = 0; i < count; i++)
		printf("%s\n", funcNames[i]);
	printf("-----------------------------------------------------\n");
	printf("END CORE DUMP\n");	
	free(funcNames);
}

int infof(const char *format, ...)
{
	textcolor(BRIGHT, WHITE, BLACK);	
	printf("INFO: ");
	va_list arg;
	int done;
	va_start (arg, format);
	done = vfprintf (stdout, format, arg);
	va_end (arg);
	printf("\033[0m");
	return done;
}

int warnf(const char *format, ...)
{
	textcolor(BRIGHT, MAGENTA, BLACK);	
	printf("WARN: ");
	va_list arg;
	int done;
	va_start (arg, format);
	done = vfprintf (stdout, format, arg);
	va_end (arg);
	printf("\033[0m");
	printf("\033[0m");
	return done;
}

int errorf(const char *format, ...)
{
	textcolor(BRIGHT, YELLOW, BLACK);	
	printf("ERROR: ");
	va_list arg;
	int done;
	va_start (arg, format);
	done = vfprintf (stdout, format, arg);
	va_end (arg);
	printf("\033[0m");
	return done;
}

int panicf(const char *format, ...)
{
	textcolor(BRIGHT, RED, BLACK);	
	printf("PANIC: ");
	va_list arg;
	int done;
	va_start (arg, format);
	done = vfprintf (stdout, format, arg);
	va_end (arg);
	print_backtrace();
	printf("\033[0m");
	exit(-1);
	return done;
}




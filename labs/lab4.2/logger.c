#include "logger.h"

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

int logDest = 0;  //0 for STDOUT, 1 for SYSLOG

int initLogger(char *logType)
{
	if((strlen(logType) == 0) || (strcmp(logType, "stdout")) == 0)
	{
		//Log to stdout
		logDest = 0;
		return 0;
	}
	else if(strcmp(logType, "syslog") == 0)
	{
		//Log to syslog
		logDest = 1;
		return 0;
	}
	panicf("Failed to initialize logger\n");
	return -1;
}

void textcolor(int attr, int fg, int bg)
{	char command[13];

	/* Command is the control command to the terminal */
	sprintf(command, "%c[%d;%d;%dm", 0x1B, attr, fg + 30, bg + 40);
	printf("%s", command);
}	

void printBacktrace(){

	void *tracePtrs[10];
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

void printBacktraceSyslog(){

	void *tracePtrs[10];
	size_t count;
	count = backtrace(tracePtrs, 10);
	char** funcNames = backtrace_symbols(tracePtrs, count);
	syslog(LOG_ERR, "BEGIN CORE DUMP\n");
	syslog(LOG_ERR, "-----------------------------------------------------\n");
	for (int i = 0; i < count; i++)
		syslog(LOG_EMERG, funcNames[i]);
	syslog(LOG_ERR, "-----------------------------------------------------\n");
	syslog(LOG_ERR, "END CORE DUMP\n");	
	free(funcNames);
}

int infof(const char *format, ...)
{
	textcolor(BRIGHT, WHITE, BLACK);
	va_list arg;
	int done;
	if(logDest == 0)
	{
		va_start (arg, format);
		printf("INFO: ");
		done = vfprintf (stdout, format, arg);
	}
	else if(logDest == 1)
	{
		va_start (arg, format);
		openlog ("Logger-INFO", LOG_CONS | LOG_PID | LOG_NDELAY, LOG_LOCAL1);
		vsyslog(LOG_INFO, format, arg);
		closelog();
	}
	va_end (arg);
	printf("\033[0m");
	return done;
}

int warnf(const char *format, ...)
{
	textcolor(BRIGHT, MAGENTA, BLACK);	
	va_list arg;
	int done;
	if(logDest == 0)
	{
		va_start (arg, format);
		printf("WARN: ");
		done = vfprintf (stdout, format, arg);
	}
	else if(logDest == 1)
	{
		va_start (arg, format);
		openlog ("Logger-WARN", LOG_CONS | LOG_PID | LOG_NDELAY, LOG_LOCAL1);
		vsyslog(LOG_WARNING, format, arg);
		closelog();
	}
	va_end (arg);
	printf("\033[0m");
	return done;
}

int errorf(const char *format, ...)
{
	textcolor(BRIGHT, YELLOW, BLACK);	
	va_list arg;
	int done;
	if(logDest == 0)
	{
		va_start (arg, format);
		printf("ERROR: ");
		done = vfprintf (stdout, format, arg);
	}
	else if(logDest == 1)
	{
		va_start (arg, format);
		openlog ("Logger-ERROR", LOG_CONS | LOG_PID | LOG_NDELAY, LOG_LOCAL1);
		vsyslog(LOG_ERR, format, arg);
		closelog();
	}
	va_end (arg);
	printf("\033[0m");
	return done;
}

int panicf(const char *format, ...)
{
	textcolor(BRIGHT, RED, BLACK);	
		va_list arg;
	int done;
	if(logDest == 0)
	{
		va_start (arg, format);
		printf("PANIC: ");
		done = vfprintf (stdout, format, arg);
		printBacktrace();
	}
	else if(logDest == 1)
	{
		va_start (arg, format);
		openlog ("Logger-PANIC", LOG_CONS | LOG_PID | LOG_NDELAY, LOG_LOCAL1);
		vsyslog(LOG_ERR, format, arg);
		printBacktraceSyslog();
		closelog();
	}
	va_end (arg);
	printf("\033[0m");
	exit(-1);
	return done;
}




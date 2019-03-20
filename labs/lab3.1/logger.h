#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <execinfo.h>
#include <string.h>
#include <syslog.h>
#include <unistd.h>


void textcolor(int attr, int fg, int bg);
int infof(const char *format, ...);
int warnf(const char *format, ...);
int errorf(const char *format, ...);
int panicf(const char *format, ...);
void printBacktrace();
void printBacktraceSyslog();
int initLogger(char *logType);

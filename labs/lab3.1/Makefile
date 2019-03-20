CC       = gcc
CFLAGS   = -Wall -rdynamic
LDFLAGS  =
OBJFILES = logger.o monitor.o
TARGET   = monitor.out

all: $(TARGET)

$(TARGET) : $(OBJFILES)
	$(CC) $(CFLAGS) -o $(TARGET) $(OBJFILES) $(LDFLAGS)

clean:
	rm -f $(OBJFILES) $(TARGET) *~
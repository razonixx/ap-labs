CC       = gcc
CFLAGS   = -rdynamic
LDFLAGS  =
OBJFILES = logger.o base64.o
TARGET   = base64.out

all: $(TARGET)

$(TARGET) : $(OBJFILES)
	$(CC) $(CFLAGS) -o $(TARGET) $(OBJFILES) $(LDFLAGS)

clean:
	rm -f $(OBJFILES) $(TARGET) *~
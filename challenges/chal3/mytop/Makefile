CC       = gcc
CFLAGS   = -Wall
LDFLAGS  =
OBJFILES = mytop.o
TARGET   = mytop.out

all: $(TARGET)

$(TARGET) : $(OBJFILES)
	$(CC) $(CFLAGS) -o $(TARGET) $(OBJFILES) $(LDFLAGS)

clean:
	rm -f $(OBJFILES) $(TARGET) *~
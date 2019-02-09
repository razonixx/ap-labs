#include <stdlib.h>
#include <fcntl.h>
#include <unistd.h>
#include <dirent.h>

#define PATH "/dev/pts/"

void myWall(int file_count, char* buffer);
char* fillBuffer(int argc, char *argv[]);
int getNumFiles(char *path);
int mystrlen(char *str);
char* stradd(char *origin, char *addition);
char* straddchar(char *origin, char addition);


int main(int argc, char *argv[])
{ 
    if(argc < 2){
        write(1, "No message to send\n", 20);
        return -1;
    }

    //Access Directory and count terminals
    int file_count = getNumFiles(PATH);

    //Buffer that will be written to the terminals
    char *buffer = fillBuffer(argc, argv);
    
    //Write buffer to each terminal
    myWall(file_count, buffer);

    return 0;
} 


//Function to write buffer to each terminal
void myWall(int file_count, char* buffer)
{
    //Send the buffer to each terminal
    for(int i = 0; i < file_count+1; i++)
    {
        int ofd = open(straddchar(PATH, (char) (i + 48)), O_WRONLY); 
        write(ofd, "\n", sizeof(char));
        write(ofd, "\n", sizeof(char));
        write(ofd, buffer, mystrlen(buffer));
        write(ofd, "\n", sizeof(char));
        write(ofd, "\n", sizeof(char));
        close(ofd);
    }
}

//Function to fill buffer with the parameters recieved from console
char* fillBuffer(int argc, char *argv[])
{
    char *buffer = "";
    //Get the complete string into the buffer
    while(argc-- > 1)
    {
        buffer = stradd(buffer, *++argv);
        buffer = stradd(buffer, " ");
    }
    return buffer;
}

//Function to count the number of files on path
int getNumFiles(char *path)
{
    int file_count = 0;
    DIR * dirp;
    struct dirent *entry;
    dirp = opendir(path);
    while ((entry = readdir(dirp)) != NULL) {
            file_count++;
    }
    closedir(dirp);
    return file_count;
}

int mystrlen(char *str)
{
    int i = 0;
    while(str[i] != '\0')
    {
        i++;
    }
    return i;
}

char* stradd(char *origin, char *addition)
{
    int lenOrigin = mystrlen(origin);
    int lenAddition = mystrlen(addition);
    char *res = (char *)malloc(lenOrigin + lenAddition + 1);
    
    int i;
    for(i = 0; i < lenOrigin; i++)
    {
        res[i] = origin[i];
    }

    for(int j = 0; j < lenAddition; j++)
    {
        res[i] = addition[j];
        i++;
    }
    res[i] = '\0';
    return res;
}

char* straddchar(char *origin, char addition)
{
    int lenOrigin = mystrlen(origin);
    char *res = (char *)malloc(lenOrigin + sizeof(char));
    
    int i;
    for(i = 0; i < lenOrigin; i++)
    {
        res[i] = origin[i];
    }

    res[i] = addition;
    return res;
}
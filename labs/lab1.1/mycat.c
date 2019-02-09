#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>

const int BUFF_SIZE = 0x6fffff;

int main(int argc, char *argv[])
{ 
    if(argc < 2){
        printf("You need to pass at least 1 file to cat");
        return -1;
    }
    char *buffer = (char *) calloc(BUFF_SIZE * (argc - 1), sizeof(char)); 
    while(--argc > 0)
    {
        int ofd = open(*++argv, O_RDONLY); 

        if (ofd < 0) { 
            perror("Error opening file"); 
            return -1;
        } 

        int rf = read(ofd, buffer, BUFF_SIZE); 

        buffer[rf] = '\0'; 
        printf("%s\n", buffer); 
        close(ofd);
    }
} 
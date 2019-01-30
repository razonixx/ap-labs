#include "stdio.h"

int strlen(char *str);
char* stradd(char *origin, char *addition);
int strfind(char *origin, char *substr);

int main(int argc, char const *argv[])
{

    if(argc != 4)
    {
        printf("You must pass 3 arguments, the string, the addition of the string and the substring\n");
        return -1;
    }

    printf("Initial length: %d\n", strlen(argv[1]));

    char *newStr = stradd(argv[1], argv[2]);

    printf("New String: %s\n", newStr);

    if(strfind(newStr, argv[3]))
    {
        printf("Substring was found: yes\n");
    }
    else
    {   
        printf("Substring was found: no\n");
    }

    return 0;
}
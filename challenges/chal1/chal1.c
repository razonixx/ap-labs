#include<stdio.h>

void ex1(char* list)
{
    int c = 0;
    int c2 = 0;
    char list2[40];
    int hasLetter = 0;

    
    for(int i = 0; i < sizeof(list)/sizeof(list[0]); i++)
    {
        hasLetter=0;
        for(int j = 0; j < sizeof(list2)/sizeof(list2[0]); j++)
        {
            if((list[i] == list2[j]) && !hasLetter)
            {
                hasLetter = 1;
            }
        }
        if(!hasLetter)
        {
            list2[c] = list[i];
        }
        c++;
    } 
    printf("%s\n", list2);
}

int main()
{
    ex1("pwwkew");
}

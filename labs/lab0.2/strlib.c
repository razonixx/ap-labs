int strlen(char *str)
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
    int lenOrigin = strlen(origin);
    int lenAddition = strlen(addition);
    char *res = (char *)malloc(lenOrigin + lenAddition);
    
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
    return res;
}

int strfind(char *origin, char *substr)
{
    int j = 0;
    for(int i = 0; i < strlen(origin); i++)
    {
        if(j != 0 && origin[i] != substr[j])
        {
            j = 0;
        }
        if(origin[i] == substr[j])    
        {
            j++;
        }
        if(j == strlen(substr))
        {
            return 1;
        }
    }
    return 0;
}
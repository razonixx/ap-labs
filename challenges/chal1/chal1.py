def ex1(string, list, l):
    if len(string) == 0:
        return l
    else:
        if string[0] not in list:
            l += 1
            list.append(string[0])
        else:
            return l
        return ex1(string[1:], list, l)

#string = "abcabcbb"
#string = "bbbbb"
#string = "qwwkew"
#string = "water"
#string = "aabbcc"
string = "abcdefghijk"
l = 0
for i in range(len(string)):
    list = []
    l2 = ex1(string[i:], list, 0)
    if l < l2:
        l = l2
print(l)

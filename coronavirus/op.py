List = [str(i) for i in range(0, 50)]
print(List)

length = len(List)
core = []
ccore = []
cccore = []
An = Bn = Cn = 0
rows = 5
columns = 4
page = 0
Next = "» Next"
Back = "« Back"
row, col = 0, 0

while True:
    if page <= 0:
        if len(List) > rows*columns - (row // 1 * columns + col):
            if row + 1 == rows and col + 1 == columns:
                core.append("» Next")
                col += 1
            else:
                core.append(List[0])
                col += 1
                List = List[1:]
        else:
            core.append(List[0])
            col += 1
            List = List[1:]
    else:
        if len(List) > rows*columns - (row // 1 * columns + col):
            if row + 1 == rows and col == 0:
                core.append("« Back")
                col += 1
            elif row + 1 == rows and col + 1 == columns:
                core.append("» Next")
                col += 1
            else:
                core.append(List[0])
                col += 1
                List = List[1:]
        else:
            if (len(List) + (row // 1 * columns + col)) // columns == row and col == 0:
                core.append("« Back")
                col += 1
                core.extend(List)
                List = List[len(List):]
            else:
                core.append(List[0])
                col += 1
                List = List[1:]
        
    if len(core) == columns or len(List) == 0:
        ccore.append(core)
        row += 1
        col = 0
        core = []

    # now = (row//(rows+1)) * rows + col
    now = row // 1 * columns + col
    if now == rows*columns or len(List) == 0:
        cccore.append(ccore)
        ccore = []
        row, col = 0, 0
        page += 1
        if len(List) == 0:
            break
print(cccore)
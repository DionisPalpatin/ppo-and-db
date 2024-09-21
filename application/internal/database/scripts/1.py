data = list()
with open(r"G:\Study\University\Database Course Project\notebook app\scripts\create.sql", mode="r", encoding="utf-8") as f:
    data = f.readlines()
    data = "".join(data)
    data = data.lower()
with open(r"G:\Study\University\Database Course Project\notebook app\scripts\create1.sql", mode="w", encoding="utf-8") as f:
    f.write(data)
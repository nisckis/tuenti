n = int(input())
data = []

for i in range(n):
    data.append(input())

d = {
    ".": "e",
    ",": "w",
    ";": "z",
    "'": "q",
    "-": "'",
    "a": "a",
    "b": "n",
    "c": "i",
    "d": "h",
    "e": "d",
    "f": "y",
    "g": "u",
    "h": "j",
    "i": "g",
    "j": "c",
    "k": "v",
    "l": "p",
    "m": "m",
    "n": "l",
    "o": "s",
    "p": "r",
    "q": "x",
    "r": "o",
    "s": ";",
    "t": "k",
    "u": "f",
    "v": ".",
    "w": ",",
    "x": "b",
    "y": "t",
    "z": "/",
}


i = 0
for case in data:
    test = ""
    for c in case:
        if d.get(c) is not None:
            test += d[c]
        else:
            test += c

    i += 1
    print(f"Case #{i}: {test}")


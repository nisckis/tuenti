from telnetlib import Telnet
from time import time
import networkx as netx
from random import randint
from os import system

moves = [
    ("2u1l", -2, -1),
    ("2u1r", -2, +1),
    ("1u2r", -1, +2),
    ("1d2r", +1, +2),
    ("2d1r", +2, +1),
    ("2d1l", +2, -1),
    ("1d2l", +1, -2),
    ("1u2l", -1, -2),
]

backs = {
    "2u1l": "2d1r",
    "2u1r": "2d1l",
    "1u2r": "1d2l",
    "1d2r": "1u2l",
    "2d1r": "2u1l",
    "2d1l": "2u1r",
    "1d2l": "1u2r",
    "1u2l": "1d2r",
}

def solve():
    tn = Telnet("52.49.91.111", 2003)
    data = tn.read_until(str.encode("--- Quick there is no time to lose! The princess is in danger ---"))
    maze = data.decode().split('\n')[:5]

    seen = ["K"]

    x0, y0 = 0, 0
    x, y = 0, 0
    base, curr, last = "K", "K", None

    while True:
        for i in range(len(maze)):
            for j in range(len(maze)):
                if maze[i][j] == "K":
                    x, y = i, j
        
        system("clear")
        print(last)
        for row in maze:
            print(row)

        stuck = True

        for m, mx, my in moves:
            if x + mx < 0 or x + mx >= len(maze) or y + my < 0 or y + my >= len(maze):
                continue

            if maze[x + mx][y + my] == "#":
                continue

            if f"{curr},{m}" in seen:
                continue
            
            stuck = False

            curr = f"{curr},{m}"
            last = m

            tn.write(str.encode(f"{m}\n"))

            if maze[x + mx][y + my] == "P":
                print("Found!")
            
            break

        if stuck:
            curr = curr[:-5]
            tn.write(str.encode(f"{backs[last]}\n"))

        maze = read_maze(tn)
        if maze is None:
            break
        


def read_maze(tn):
    maze = []

    while True:
        try:
            d = tn.read_until(str.encode("\n")).decode().replace("\n", "")
        except Exception as e:
            print(e)
            print(d)
            print(maze)
            return None

        data = d

        if len(data) > 0:
            maze.append(data)

        if len(data) == 0 and len(maze) > 0:
            break
    
    return maze

if __name__ == "__main__":
    solve()
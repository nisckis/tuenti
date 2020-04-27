from telnetlib import Telnet
from time import time
import networkx as netx
from random import randint, choice
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

    t = time()
    dx = [2, 2, -2, -2, 1, 1, -1, -1] 
    dy = [1, -1, 1, -1, 2, -2, 2, -2]

    its = 0

    G = netx.DiGraph()
    G.add_node("K")
    G.nodes["K"]["visited"] = False
    G.nodes["K"]["pos"] = (2, 2)

    s = ["K"]
    can_end, final, ppos = False, None, None
    last = None

    data = tn.read_until(str.encode("--- Quick there is no time to lose! The princess is in danger ---"))
    maze = data.decode().split('\n')[:5]

    while len(s) > 0:
        px, py, = 2, 2

        rv = randint(0, len(s)-1)
        v = s[rv]
        del s[rv]

        # v = s[0]
        # s = s[1:]
        
        if len(s) > 0:
            path = netx.dijkstra_path(G, last, v)

            for i in range(1, len(path)):
                move = G.edges[path[i-1], path[i]]["move"]
                tn.write(str.encode(f"{move}\n"))

            mm = []

            while True:
                try:
                    d = tn.read_until(str.encode("\n")).decode().replace("\n", "")
                except Exception as e:
                    print(f"\n ------\nException {e}\nDATA\n{d}\nLAST_MAZE\n{maze}\nMAZE\n{mm}\n ------\n")
                    return

                data = d

                if len(data) > 0 and data[0] != "-":
                    mm.append(data)

                if len(data) == 0 and len(mm) > 0:
                    break

            maze = mm

        l = len(maze)

        parcial = time() - t
        print(f"it({its:5d}, {parcial:5.3f} s.)  {len(s)}") # root is {v} and last {last}")
        its += 1

        if not G.nodes[v].get("visited"):
            G.nodes[v]["visited"] = True

            for move, x, y in moves:
                nx, ny = px + x, py + y

                if maze[nx][ny] == "#" or nx < 0 or nx >= len(maze) or ny < 0 or ny >= len(maze):
                    continue

                if maze[nx][ny] == "P":
                    print(maze)

                    can_end = True
                    G.add_edge(v, f"{v},P", move=move, back=backs[move])
                    G.add_edge(f"{v},P", v, move=backs[move], back=move)
                    final = move
                    G.nodes[f"{v},P"]["pos"] = (nx, ny)
                    G.nodes[f"{v},P"]["visited"] = False
                    ppos = f"{v},P"
                    print("found the P!")
                    break

                G.add_edge(v, f"{v},{x},{y}", move=move, back=backs[move])
                G.add_edge(f"{v},{x},{y}", v, move=backs[move], back=move)
                G.nodes[f"{v},{x},{y}"]["pos"] = (nx, ny)
                G.nodes[f"{v},{x},{y}"]["visited"] = False
                s.append(f"{v},{x},{y}")

            if can_end:
                break
        
        last = v

    print(can_end)
    print(len(G.nodes))
    print(len(G.edges))

    t = time() - t
    print(f"elapsed: {t}")

    if not can_end:
        print("wtf")
        return False

    print(f"final move is {final}")
    tn.write(str.encode(f"{final}\n"))
    maze = []
    while True:
        try:
            d = tn.read_until(str.encode("\n")).decode().replace("\n", "")
        except Exception as e:
            print(f"\n ------\nException {e}\nDATA\n{d}\nLAST_MAZE\n{maze}\n ------\n")
            return False

        print(d)
        data = d
        if len(data) > 0:
            i += 1
            maze.append(data)
        if len(data) == 0 and len(maze) > 0:
            break
    print(maze)

    path = netx.dijkstra_path(G, "K", ppos)
    print(f"the path is\n{path}")

    tn = Telnet("52.49.91.111", 2003)
    data = tn.read_until(str.encode("--- Quick there is no time to lose! The princess is in danger ---"))

    for i in range(1, len(path)):
        move = G.edges[path[i-1], path[i]]["move"]
        tn.write(str.encode(f"{move}\n"))

        maze = []

        while True:
            try:
                d = tn.read_until(str.encode("\n")).decode().replace("\n", "")
            except Exception as e:
                print(f"Raised exp: {e}\ndata\n{d}\nmaze\n{maze}")
                return False

            print(d)
            data = d

            if len(data) > 0:
                i += 1
                maze.append(data)

            if len(data) == 0 and len(maze) > 0:
                break

        print(maze)

    return True

if __name__ == "__main__":
    while True:
        try:
            if solve():
                break
        except:
            continue
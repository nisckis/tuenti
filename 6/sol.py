from telnetlib import Telnet
import socket
from time import time
import networkx as netx
from collections import defaultdict
from random import randint, choice
from math import sqrt
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

moves_obj = {
    "2u1l": [-2, -1],
    "2u1r": [-2, +1],
    "1u2r": [-1, +2],
    "1d2r": [+1, +2],
    "2d1r": [+2, +1],
    "2d1l": [+2, -1],
    "1d2l": [+1, -2],
    "1u2l": [-1, -2],
}

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
    sock = socket.socket()
    sock.connect(("52.49.91.111", 2003))

    data = sock.recv(128)
    maze = data.decode().split('\n')[:5]

    # tn = Telnet("52.49.91.111", 2003)

    t = time()
    its = 0
    expanded = 0

    N, M = 130, 160

    # G = netx.DiGraph()
    G = netx.DiGraph()
    G.add_nodes_from((i, j) for i in range(N) for j in range(M))

    G.add_node("K")
    G.nodes["K"]["visited"] = False
    G.nodes["K"]["pos"] = (0, 0)

    seen = defaultdict(lambda: defaultdict(lambda: False))
    # actual = [0, 0]
    actual = [5 * N // 6, 2 * M // 4]

    m = []
    for _ in range(N):
        x = []
        for _ in range(M):
            x.append("#")
        m.append(x)

    m_pos = [actual[0], actual[1]]
    s = [(actual[0], actual[1])]

    can_end, final, ppos = False, None, None
    last = None

    # data = tn.read_until(str.encode(
    #     "--- Quick there is no time to lose! The princess is in danger ---"))
    # maze = data.decode().split('\n')[:5]

    while len(s) > 0:
        px, py, = 2, 2
        for i in range(len(maze)):
            for j in range(len(maze[i])):
                if maze[i][j] == "K":
                    px, py = i, j

        # rv = randint(0, len(s)-1)
        # v = s[rv]
        # del s[rv]

        # if last is not None:
        #     s = sorted(s, key=lambda x: sqrt(x[0] * last[0] + x[1] * last[1]))

        v = s[0]
        s = s[1:]

        # if G.nodes[v].get("visited"):
        #     continue

        # v = s[-1]
        # s = s[:-1]

        t_rep = time()
        t_rep2 = 0

        # Move the K to the v position
        if last is not None:
            # path = netx.dijkstra_path(G, last, v)
            path = netx.shortest_path(G, last, v)

            for i in range(1, len(path)):
                move = G.edges[path[i-1], path[i]]["move"]

                actual[0] += moves_obj[move][0]
                actual[1] += moves_obj[move][1]

                m_pos[0] += moves_obj[move][0]
                m_pos[1] += moves_obj[move][1]

                t_tmp = time()
                sock.send(str.encode(f"{move}\n"))
                data = sock.recv(128)
                t_rep2 += (time() - t_tmp)

            data = data.decode().split('\n')
            
            if len(data) == 2:
                print(f"this {move} has some problems from {last} to {v}")
                print(f"path is {path}")
                for d in data:
                    print(d)
                sock.close()
                return False

            if data[0][0] == "-":
                print(f"this {move} has some problems from {last} to {v}")
                print(f"path is {path}")
                for d in data:
                    print(d)
                maze = data[1:6]
            else:
                maze = data[:5]

            if len(maze) != 5:
                raise Exception("wtf")

            for i in range(len(maze)):
                for j in range(len(maze[i])):
                    m[m_pos[0]-2+i][m_pos[1]-2+j] = maze[i][j]

                

            # tn.write(str.encode(f"{move}\n"))
            # mm = []
            # while True:
            #     try:
            #         d = tn.read_until(str.encode(
            #             "\n")).decode().replace("\n", "")
            #     except Exception as e:
            #         print(
            #             f"\n------\nException {e}\nDATA\n{d}\n"
            #             f"LAST_MAZE\n{maze}\nMAZE\n{mm}\n------\n"
            #             )
            #         return
            #     data = d
            #     if len(data) > 0 and data[0] != "-":
            #         mm.append(data)
            #     if len(data) == 0 and len(mm) > 0:
            #         break
            # maze = mm

        t_rep = time() - t_rep
        t_main = time()

        # if not G.nodes[v].get("visited") and not seen[str(actual[0])][str(actual[1])]:
        if not seen[str(actual[0])][str(actual[1])]:
            G.nodes[v]["visited"] = True

            seen[str(actual[0])][str(actual[1])] = True

            # print(actual, seen_x, seen_y)

            for move, x, y in moves:
                nx, ny = px + x, py + y

                # Check if move is valid
                if maze[nx][ny] == "#":
                    continue

                # if nx < 0 or nx >= len(maze) or ny < 0 or ny >= len(maze):
                #     continue

                if maze[nx][ny] == "P":
                    print(data)
                    print(maze)

                    # can_end = True
                    # G.add_edge(v, f"{v},P", move=move, back=backs[move])
                    # G.add_edge(f"{v},P", v, move=backs[move], back=move)
                    # final = move
                    # G.nodes[f"{v},P"]["pos"] = (nx, ny)
                    # G.nodes[f"{v},P"]["visited"] = False
                    # ppos = f"{v},P"

                    can_end = True
                    G.add_edge((actual[0], actual[1]), (actual[0]+x, actual[1]+y), move=move)
                    G.add_edge((actual[0]+x, actual[1]+y), (actual[0], actual[1]), move=backs[move])
                    G.nodes[(actual[0]+x, actual[1]+y)]["pos"] = (nx, ny)
                    G.nodes[(actual[0]+x, actual[1]+y)]["visited"] = False
                    s.append((actual[0]+x, actual[1]+y))
                    ppos = (actual[0]+x, actual[1]+y)
                    expanded += 1

                    print("found the P!")
                    break

                # Check if adding this move moves you to a seen position
                if seen[str(actual[0]+x)][str(actual[1]+y)]:
                    continue

                # G.add_edge(v, f"{v},{move}", move=move, back=backs[move])
                # G.add_edge(f"{v},{move}", v, move=backs[move], back=move)
                # G.nodes[f"{v},{move}"]["pos"] = (nx, ny)
                # G.nodes[f"{v},{move}"]["visited"] = False
                # s.append(f"{v},{move}")
                # expanded += 1

                G.add_edge((actual[0], actual[1]), (actual[0]+x, actual[1]+y), move=move)
                G.add_edge((actual[0]+x, actual[1]+y), (actual[0], actual[1]), move=backs[move])
                G.nodes[(actual[0]+x, actual[1]+y)]["pos"] = (nx, ny)
                G.nodes[(actual[0]+x, actual[1]+y)]["visited"] = False
                s.append((actual[0]+x, actual[1]+y))
                expanded += 1

                # print(f"added from {v} to {v},{move}")

            if can_end:
                break

        

        # system("cls")
        
        t_main = time() - t_main
        parcial = time() - t
        
        # for row in m:
        #     for c in row:
        #         print(c, end="")
        #     print()
        its += 1
        print(f"#{its:5d} @ {parcial:f} s. -> {(time() - t)/ its:f} s./it. ({t_rep2:f} / {t_rep:f} s.) ({t_main:f} s.) - {expanded} / {len(s)} - ({actual[0]}, {actual[1]}) {v}")

        last = v

    for row in m:
        for c in row:
            print(c, end="")
        print()

    print(can_end)
    print(len(G.nodes))
    print(len(G.edges))

    t = time() - t
    print(f"elapsed: {t}")

    sock.close()

    # if not can_end:
    #     print("wtf")
    #     return False

    sock = socket.socket()
    sock.connect(("52.49.91.111", 2003))

    data = sock.recv(128)
    maze = data.decode().split('\n')[:5]

    if ppos is not None:
        path = netx.dijkstra_path(G, (20, 40), ppos)
    else:
        path = netx.dijkstra_path(G, (20, 40), (21, 43))

    print(f"the path is\n{path}")

    # tn = Telnet("52.49.91.111", 2003)
    # data = tn.read_until(str.encode(
    #     "--- Quick there is no time to lose! The princess is in danger ---"))

    for i in range(1, len(path)):
        move = G.edges[path[i-1], path[i]]["move"]

        sock.send(str.encode(f"{move}\n"))
        maze = sock.recv(128).decode().split('\n')[:5]

        # tn.write(str.encode(f"{move}\n"))
        # maze = []
        # while True:
        #     try:
        #         d = tn.read_until(str.encode("\n")).decode().replace("\n", "")
        #     except Exception as e:
        #         print(f"Raised exp: {e}\ndata\n{d}\nmaze\n{maze}")
        #         sock.close()
        #         return False
        #     print(d)
        #     data = d
        #     if len(data) > 0:
        #         i += 1
        #         maze.append(data)
        #     if len(data) == 0 and len(maze) > 0:
        #         break

        print(maze)

    sock.close()
    return True


if __name__ == "__main__":
    solve()

    # while True:
    #     try:
    #         if solve():
    #             break
    #     except:
    #         continue

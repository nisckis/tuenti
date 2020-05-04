from PIL import Image
import numpy as np

original = Image.open("original.jpg")

target = Image.open("original.png")
original = original.resize(target.size)

print(target.size)
print(original.size)


org = np.asarray(original)
tar = np.asarray(target)

size = target.size

count = 0

for x in range(size[1]):
    for y in range(size[0]):
        if tar[x][y].all() != org[x][y].all():
            count += 1

print(count, size[0] * size[1])
print(size[0] * size[1] - count)

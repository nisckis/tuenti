def gcd(a, b):
    if a == 0:
        return b, 0, 1

    g, x1, y1 = gcd(b % a, a)

    x = y1 - (b // a) * x1
    y = x1

    return g, x, y


file_pairs = [
    ("testdata/plaintexts/test1.txt", "testdata/ciphered/test1.txt"),
    ("testdata/plaintexts/test2.txt", "testdata/ciphered/test2.txt")
]

e = 13 # 65537
test = []

print("-" * 60)


for plaintext, ciphered in file_pairs:
    print(f"{plaintext} -> {ciphered}")

    plain_b = open(plaintext, "rb").read()
    plain_lines = open(plaintext, "rb").readlines()
    plain = [c for line in plain_lines for c in line]

    cipher_b = open(ciphered, "rb").read()
    cipher_lines = open(ciphered, "rb").readlines()
    cipher = [c for line in cipher_lines for c in line]

    print(plain)
    print(cipher)

    test.append([plain[i]**e - cipher[i] for i in range(len(plain))])

    print("-" * 60)

for a, b in zip(test[0], test[1]):
    print(a, b, gcd(a, b)[0])

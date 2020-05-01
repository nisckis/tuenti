from Crypto.Util.Padding import pad, unpad
from Crypto.Cipher import AES
from Crypto.Random import get_random_bytes

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

e = 3 # 65537
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
    print(int(plain_b))

    test.append([plain[i]**e - cipher[i] for i in range(len(plain))])

    print("-" * 60)

data = b'First test'
key = get_random_bytes(32)
iv = get_random_bytes(16)

cipher1 = AES.new(key, AES.MODE_CBC, iv)
ct = [c for c in cipher1.encrypt(pad(data, 128))]

print(ct)

# for a, b in zip(test[0], test[1]):
#     print(a, b, gcd(a, b)[0])

# from string import ascii_lowercase as a_lower

# for c in ["6D", "65", "73", "73", "61", "67", "65"]:
#     print(chr(int(c, 16)))


cipher = "3633363A33353B393038383C363236333635313A353336"
target = "514;248;980;347;145;332"

# test = []

cipher_split = [cipher[i:i+2] for i in range(0, len(cipher), 2)]
target_hex = [c for c in target]
target_split = [c for c in target]

# print(" ".join(cipher_split))
# print(" ".join([" " + chr(int(c, 16)) for c in cipher_split]))
# print(" ".join([" " + c for c in target]))

cipher_dec = [int(c, 16) for c in cipher_split]
target_dec = [ord(c) for c in target_split]

# print(cipher_dec)
# print(target_dec)

xor = [
    [ n for n in [i for i in range(11)] if a ^ n == b]
    for a, b in zip(cipher_dec, target_dec)
]

wtf = [str(x[0]) for x in xor]
key = wtf[::-1]

print("key", "".join(key))

juj = "3A3A333A333137393D39313C3C3634333431353A37363D"
juj_split = [juj[i:i+2] for i in range(0, len(juj), 2)]

print(juj_split)
print(wtf)

data = [
    int(j, 16) ^ int(w)
    for j, w in zip(juj_split, wtf)
]

print("".join(chr(x) for x in data))


def encrypt(key, msg):
    cip = ""

    for i in range(len(msg)):
        c = msg[i]
        asc_chr = ord(c)
        key_char = key[len(key) - 1 - i]
        
        print(c, asc_chr, key_char)

        crpt_chr = asc_chr ^ ord(key_char)
        print("XOR", asc_chr, key_char, crpt_chr) 

        hx_crpt_chr = hex(crpt_chr)
        cip += " " + hx_crpt_chr

    return cip


# print(encrypt("somekey", "tarasyarema"))

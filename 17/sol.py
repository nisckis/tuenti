import zlib
import binascii

IDAT = 'A0A5E2A091E2A09DE2A09EE2A08AE2A089E2A093E2A081E2A087E2A087E2A091E2A09DE2A09BE2A091E2A0BCE2A081E2A0BCE2A09A'
data = bytes.fromhex(IDAT)
result = binascii.hexlify(zlib.decompress(data))


# print(result) 
print(result)
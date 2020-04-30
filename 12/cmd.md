# Generate keys
openssl genrsa -out key.pem 1024
openssl rsa -in key.pem -pubout > key.pub

# Encrypt file with pub key
openssl rsautl -encrypt -pubin -inkey key.pub -in text -out ciphertext


# Read cipher as hex
xxd -ps -l 128 cipher.txt

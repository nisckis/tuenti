from paramiko import SSHClient, AutoAddPolicy
from emoji import EMOJI_UNICODE as emos
from emoji import emojize

ssh = SSHClient()
ssh.set_missing_host_key_policy(AutoAddPolicy())

ssh.connect(
  "52.49.91.111",
  port=22000,
  username="castle",
  password="castle",
  look_for_keys=False,
)

chan = ssh.invoke_shell()

for k in emos:
  print(f"Trying emoji {k} -> {emos[k]} {emojize(k)}\n")
  
  try:
    while not chan.recv_ready():
      pass

    b1 = chan.recv(1024).decode()

    while not chan.send_ready():
      pass

    chan.send(f"{emojize(k)}\n")

    b2 = chan.recv(1024).decode()

    while not chan.recv_ready():
      pass

    b3 = chan.recv(1024).decode()
  
    while not chan.send_ready():
      pass

    chan.send(f"{emojize(k)}\n")

    b4 = chan.recv(1024).decode()
    
    if "not found" in b1 and "not found" in b2 and "not found" in b3 and "not found" in b4:
      raise Exception("BAD!")

    print(f"OK! emoji {k} -> {emos[k]}\n")
    print(1, b1)
    print(2, b2)
    print(3, b3)
    print(4, b4)

    
  except Exception as e:
    # pass
    print(f"got and exception: {e}")

ssh.close()



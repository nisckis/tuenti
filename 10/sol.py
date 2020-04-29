from paramiko import SSHClient, AutoAddPolicy
import emoji

emos = emoji.EMOJI_UNICODE

with SSHClient() as ssh:
    ssh.set_missing_host_key_policy(AutoAddPolicy())
    ssh.connect(
        "52.49.91.111",
        port=22000,
        username="castle",
        password="castle",
        look_for_keys=False,
    )

    chan = ssh.invoke_shell()

    while True:
        while not chan.recv_ready():
            pass
        
        b = ""
        b = chan.recv(16384).decode()

        print("-" * 10)
        print(b)      
        print("-" * 10)
          
        cmd = input("> ")

        if cmd[:2] == "ll":
            cmd = "🔦 -al"

        if cmd[:2] == "ls":
            cmd = "🔦 " + cmd[3:]

        if cmd[:2] == "cd":
            cmd = "🚶 " + cmd[3:]

        if cmd[:2] == "gd":
            cmd = "🚶 🚪" + cmd[3:]

        if cmd[:4] == "talk":
            cmd = "💬 " + cmd[5:]

        if cmd[:4] == "grab":
            cmd = "✋ " + cmd[5:]

        if cmd[:4] == "tell":
            cmd = "🗣️ " + cmd[5:]

        msg = cmd + "\n"

        # print(f"recieved '{cmd}' and sending '{msg}'")

        while not chan.send_ready():
            pass
        chan.send(msg)

    ssh.close()

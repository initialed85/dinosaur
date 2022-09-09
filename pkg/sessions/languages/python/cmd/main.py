import os
import socket
import time
from threading import Thread

PORT = 13337


def receive_callback(sock: socket.socket, data, addr, local_ip):
    if addr == (local_ip, PORT):
        return

    print(f"{addr[0]}:{addr[1]}\t{data.decode('utf-8')}")


def receive_loop(sock, local_ip):
    while 1:
        try:
            data, addr = sock.recvfrom(65507)
        except socket.timeout:
            continue

        receive_callback(sock, data, addr, local_ip)


def main():
    hostname = os.getenv("HOSTNAME")
    local_ip = os.getenv("LOCAL_IP")
    broadcast_ip = os.getenv("BROADCAST_IP")

    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEPORT, 1)
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
    sock.bind(("0.0.0.0", PORT))
    sock.settimeout(1)

    thread = Thread(
        target=receive_loop,
        args=(
            sock,
            local_ip,
        ),
    )
    thread.start()

    data = f"Hello world from Python @ {hostname}".encode("utf-8")

    while 1:
        sock.sendto(data, (broadcast_ip, PORT))
        time.sleep(1)


if __name__ == "__main__":
    main()

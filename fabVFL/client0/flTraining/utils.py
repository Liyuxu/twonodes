import pickle
import socket
import struct
import sys

def recvall(sock):
    BUFF_SIZE = 4096  # 4 KiB
    data = b''
    while True:
        part = sock.recv(BUFF_SIZE)
        data += part
        if len(part) < BUFF_SIZE:
            # either 0 or end of data
            break
    return data

def send_msg(sock, msg):
    # bytes("LOCAL MODEL CAN LOAD"+"#"+str(num_epoch), encoding="UTF-8")
    msg_pickle = bytes(msg, encoding="UTF-8")
    sock.sendall(msg_pickle)
    print(msg, 'sent to', sock.getpeername())


def recv_msg(sock, expect_msg_type=None):
    msg = recvall(sock)
    str_msg = msg.decode('UTF-8')
    # 解析信号，使用”#“分割
    str_split = str_msg.split('#')
    if (expect_msg_type is not None) and (str_split[0] != expect_msg_type):
        print(">> msg = recvall(sock):", msg, "--len:", len(msg))
        if len(msg) == 0:
            return "MISSION_END"
        raise Exception("Expected " + expect_msg_type + " but received " + str_split[0])
    return str_split

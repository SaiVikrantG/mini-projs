import socket
import sys


def send_request(site):
    con = socket.socket()

    try:
        con.connect((site, 80))
        print(f"Succesfully connected to {site}")

        message = f"GET / HTTP/1.1\r\nHost: {site}\r\nConnection: close\r\n\r\n"
        message = message.encode("ISO-8859-1")

        con.sendall(message)

        while True:
            d = con.recv(4096)
            print(d.decode("ISO-8859-1"))

            if len(d) == 0:
                break

        con.close()
    except Exception as e:
        print(f"Exception occured: {str(e)}")


def main():
    site = sys.argv[1]

    print(site)
    send_request(site)


main()

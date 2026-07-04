import sys
import socket

# How many bytes is the word length?
WORD_LEN_SIZE = 2

BUFFER_SIZE = 2


def usage():
    print("usage: wordclient.py server port", file=sys.stderr)


packet_buffer = b""


def get_next_word_packet(s: socket.socket):
    """
    Return the next word packet from the stream.

    The word packet consists of the encoded word length followed by the
    UTF-8-encoded word.

    Returns None if there are no more words, i.e. the server has hung
    up.
    """
    global packet_buffer
    global BUFFER_SIZE

    while True:
        data = s.recv(BUFFER_SIZE)

        packet_buffer += data

        if packet_buffer == b"":
            return

        if len(packet_buffer) < WORD_LEN_SIZE:
            continue

        length_val = packet_buffer[:WORD_LEN_SIZE]
        length = int.from_bytes(length_val, byteorder="big")

        if len(packet_buffer) < WORD_LEN_SIZE + length:
            continue

        enc_word = packet_buffer[: (length + WORD_LEN_SIZE)]
        packet_buffer = packet_buffer[(length + WORD_LEN_SIZE) :]

        return enc_word


def extract_word(word_packet):
    """
    Extract a word from a word packet.

    word_packet: a word packet consisting of the encoded word length
    followed by the UTF-8 word.

    Returns the word decoded as a string.
    """

    enc_word = word_packet[2:]
    return enc_word.decode()


# Do not modify:


def main(argv):
    try:
        host = argv[1]
        port = int(argv[2])
    except Exception:
        usage()
        return 1

    s = socket.socket()
    s.connect((host, port))

    print("Getting words:")

    while True:
        word_packet = get_next_word_packet(s)

        if word_packet is None:
            break

        word = extract_word(word_packet)

        print(f"    {word}")

    s.close()


if __name__ == "__main__":
    sys.exit(main(sys.argv))

import configparser
import copy
import socket
import time

from autologging import logged

from fishbeats.proto.note_pb2 import Note


@logged
class Client:
    def __init__(self, config: configparser.ConfigParser):
        self.config = copy.copy(config)

    def start(self):
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
            for _ in range(self.config.getint("client", "connection_retries")):
                try:
                    sock.connect((self.config.get("client", "host"), self.config.getint("client", "port")))
                    self.__log.info("Connected to %s", sock)
                    for i in range(20, 61):
                        note = Note()
                        note.track = 0
                        note.num = i
                        self.__log.debug(str(note).strip())
                        sock.send(note.SerializeToString() + b'\0')
                        time.sleep(0.2)
                    break
                except ConnectionRefusedError as ex:
                    self.__log.warn(ex)
                    time.sleep(self.config.getint("client", "polling_delay"))


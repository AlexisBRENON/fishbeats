import configparser
import socketserver

from autologging import logged

from fishbeats.engine import Engine
from fishbeats.proto.note_pb2 import Note


class Server(socketserver.ThreadingTCPServer):

    def __init__(self, config: configparser.ConfigParser, engine: Engine):
        @logged
        class RequestHandler(socketserver.StreamRequestHandler):
            def notes(self) -> Note:
                while True:
                    data = self.request.recv(4096).strip()
                    if len(data) == 0:
                        break
                    for note_data in data.split(b'\0'):
                        if len(note_data) == 0:
                            continue
                        note = Note()
                        note.ParseFromString(note_data)
                        yield note

            def handle(self):
                for note in self.notes():
                    self.__log.debug("New note received: %s", str(note).strip())
                    engine.input_queue.put([note])

        super().__init__(
            (config.get("server", "host"), config.getint("server", "port")),
            RequestHandler)

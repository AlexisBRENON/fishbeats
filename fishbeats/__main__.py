#! /usr/bin/env python3
import configparser
import logging
import threading

from autologging import logged

from fishbeats.client import Client
from fishbeats.engine import Engine
from fishbeats.server import Server
from fishbeats.ui import GuiDemoClient


@logged
def serve():
    serve._log.info("Instantiate engine")
    engine_config = configparser.ConfigParser()
    engine_config.read("./share/engine.ini")
    engine = Engine(engine_config)

    serve._log.info("Instantiate server")
    server_config = configparser.ConfigParser()
    server_config.read("./share/server.ini")
    with Server(server_config, engine) as server:
        serve._log.info("Start server")
        server.serve_forever()


def request():
    client_config = configparser.ConfigParser()
    client_config.read("./share/client.ini")
    client = Client(client_config)
    client.start()


def demo():
    threading.Thread(target=serve, daemon=True).start()
    request()


def gui():
    gui = GuiDemoClient(None)
    gui.start()


if __name__ == "__main__":
    logging.basicConfig(level="DEBUG")
    gui()

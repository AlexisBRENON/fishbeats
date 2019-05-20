import configparser
import copy
import datetime
import queue
import threading
import time
from contextlib import contextmanager
from typing import Iterable

import fluidsynth
from autologging import logged

from fishbeats.proto.note_pb2 import Note


@logged
class Engine:

    def __init__(self, config: configparser.ConfigParser):
        self.config = copy.copy(config)

        self.synth = fluidsynth.Synth()
        self.synth.program_select(0, self.synth.sfload(config.get("engine", "soundfont_path")), 0, 0)
        self.events_delay: datetime.timedelta = datetime.timedelta(seconds=60) / self.config.getint("engine", "bpm") / 16
        self.main_thread = None
        self.stop_event = threading.Event()
        self.input_queue: queue.Queue[Note] = queue.LifoQueue()

    def __del__(self):
        self.synth.delete()

    def _event_loop(self):
        self.synth.start(driver=self.config.get("engine", "driver"))
        tracks = {}
        while not self.stop_event.is_set():
            loop_start = time.monotonic()
            try:
                notes: Iterable[Note] = self.input_queue.get_nowait()
                for note in notes:
                    if (note.track in tracks and
                            tracks.get(note.track) != note.num):
                        self.synth.noteoff(note.track, tracks[note.track])
                    self.synth.noteon(note.track, note.num, 30)
                    tracks[note.track] = note.num
                    self.input_queue.task_done()
                while not self.input_queue.empty():
                    self.input_queue.get_nowait()
                    self.input_queue.task_done()
            except queue.Empty:
                pass
            delay = self.events_delay.total_seconds() - (time.monotonic() - loop_start)
            self.stop_event.wait(delay)

    @contextmanager
    def start(self):
        if not self.main_thread:
            self.__log.debug("Launching engine loop")
            self.stop_event.clear()
            self.main_thread = threading.Thread(target=self._event_loop)
            self.main_thread.start()
            try:
                yield self
            except BaseException:
                raise
            finally:
                self.__log.debug("Exiting engine loop")
                self._stop()
                self.main_thread.join()

    def _stop(self):
        self.__log.info("Stopping engine")
        self.stop_event.set()

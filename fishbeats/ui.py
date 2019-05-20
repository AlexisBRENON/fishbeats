import configparser
import copy
import random
import threading
from tkinter import *
from typing import Tuple, List

from fishbeats.client import Client


class Aquarium(Frame):
    def __init__(self, *args, **kwargs):
        super(Aquarium, self).__init__(*args, width=400, height=400, **kwargs)
        self.fishes: List[Widget] = []

    def update_fish(self, fish_index: int, coords: Tuple[float, float]):
        while fish_index > len(self.fishes):
            self.update_fish(len(self.fishes), (0.5, 0.5))
        if fish_index == len(self.fishes):
            self.fishes.append(Label(self, text="O"))
        self.fishes[fish_index].place(relx=coords[0], rely=coords[1], anchor="center")


class GuiDemoClient(Tk, Client):
    def __init__(self, config: configparser.ConfigParser):
        super(GuiDemoClient, self).__init__()
        self.config = copy.copy(config)
        self.fish_list: List[Tuple[float, float]] = []
        self.stop_event = threading.Event()

        self.main_frame = Frame()
        self.main_frame.pack()

        self.aquarium = Aquarium(self.main_frame)
        self.aquarium.grid(column=1, row=1, rowspan=20)

        num_fishes_frame = Frame(master=self.main_frame)
        num_fishes_frame.grid(column=2, row=1)
        Label(master=num_fishes_frame, text="Number of fishes").pack(side="left")

        spinbox = Spinbox(master=num_fishes_frame, cnf={
            'from': 0, 'to': 10,
        })

        def spinbox_cb():
            num_fishes = int(spinbox.get())
            if num_fishes < len(self.fish_list):
                self.fish_list = self.fish_list[:num_fishes]
            else:
                self.fish_list += [(random.random(), random.random()) for _ in range(num_fishes - len(self.fish_list))]
        spinbox.configure(command=spinbox_cb)
        spinbox.pack(side="left")

    def move_ticker(self):
        while not self.stop_event.is_set():
            self.fish_list = [
                (
                    max(0, min(1, f[0] + ((random.random() * 2) - 1) / 50)),
                    max(0, min(1, f[1] + ((random.random() * 2) - 1) / 50))
                )
                for f in self.fish_list
            ]
            for i, fish in enumerate(self.fish_list):
                self.aquarium.update_fish(i, fish)
            self.stop_event.wait(0.05)

    def start(self):
        move_thread = threading.Thread(target=self.move_ticker)
        move_thread.start()
        self.mainloop()
        self.stop_event.set()
        move_thread.join()



import importlib
import sys
import os

def spawn(d, name, ws):
    path = "amebas.{}.Ameba".format(d)
    getattr(importlib.import_module(path), "Ameba")(name, ws)

ws = "ws://localhost:3000/ws"

amebas = [("Volvox", "Volvox")]

n = int(sys.argv[1]) - len(amebas)

for i in range(n):
    amebas.append(("Bot", "Bot" + str(i)))

for (d, name) in amebas:
    pid = os.fork()

    if pid == 0:
        spawn(d, name, ws)
        sys.exit()
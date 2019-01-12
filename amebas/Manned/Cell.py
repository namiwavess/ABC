import websocket

class Cell:
  def __init__(self, name, url):
    self.number = -1
    self.name = name
    self.x = 0
    self.y = 0
    self.mass = 0
    self.width = 0
    self.height = 0
    self.food = [(-1, -1)] * 100
    self.players = [(-1, -1, -1, -1)]
    self.mousex = 0
    self.mousey = 0

    websocket.enableTrace(True)
    ws = websocket.WebSocketApp(url, 
      on_message = self.on_message,
      on_error = self.on_error,
      on_close = self.on_close)
    ws.run_forever()

  def on_message(self, ws, message):
    cmd = message.split(" ")

    if cmd[0] == "1":
      self.number = int(cmd[1])
      for (i, x, y) in zip(*[map(int, cmd[2].split(",")[0:-1])]*3):
        self.food[i] = (x, y)

      self.width = int(cmd[3])
      self.height = int(cmd[4])
      
      self.initXY()
      ws.send("2 {} {} {} {}".format(cmd[1], self.name, self.x, self.y))
      
    elif cmd[0] == "4":
      for (i, x, y) in zip(*[map(int, cmd[1].split(",")[0:-1])]*3):
        self.food[i] = (x, y)

      self.players = []

      for (i, x, y, r) in zip(*[iter(cmd[2].split(",")[0:-1])]*4):
        idx = int(i)
        self.players.append((idx, int(x), int(y), float(r)))
        if idx == self.number:
          self.x = self.players[-1][1]
          self.y = self.players[-1][2]
          self.mass = self.players[-1][3]

      self.mousex = int(cmd[3].split(",")[0])
      self.mousey = int(cmd[3].split(",")[1])

      dx, dy = self.play()
      ws.send("3 {} {} {}".format(self.number, dx, dy))
    else:
      print(message)


  def on_error(self, ws, error):
    #print("debug: called on_error")
    #print(error)
    pass

  def on_close(self, ws):
    #print("killed {}".format(self.number))
    print("\u001B[6B")
    exit()

  def initXY(self):
    pass

  def play(self):
    pass

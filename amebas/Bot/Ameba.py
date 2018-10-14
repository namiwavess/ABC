from .Cell import Cell
import random
import math

class Ameba(Cell):
  def __init__(self, name, url):
    """
    self.x: ameba's position x
    self.y: ameba's position y
    self.mass: ameba's radius
    self.width = 0: world width
    self.height = 0: world height
    self.food: food positions
    self.players:[(number, x, y, mass)]
    """
    super().__init__(name, url)

  def initXY(self):
    self.x = random.randint(0, self.width)
    self.y = random.randint(0, self.height)

  def play(self):
    angle = random.randint(0, 360)
    
    speed = 3
    vx = speed * math.cos(angle)
    vy = speed * math.sin(angle)

    return int(vx), int(vy)
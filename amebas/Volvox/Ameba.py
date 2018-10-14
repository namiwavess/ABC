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

  def cal_distance(self, x, y):
    return math.sqrt(math.pow(x - self.x, 2) + math.pow(y - self.y, 2))

  def cal_vector(self, dX, dY):
    angle = math.atan2(dY - self.y, dX - self.x)

    if angle < 0:
      angle += 2 * math.pi

    speed = 3
    vx = speed * math.cos(angle)
    vy = speed * math.sin(angle)

    return int(vx), int(vy)

  def initXY(self):
    self.x = random.randint(0, self.width)
    self.y = random.randint(0, self.height)

  def play(self):
    food = sorted(self.food, key=lambda x:self.cal_distance(x[0], x[1]))
    return self.cal_vector(food[0][0], food[0][1])
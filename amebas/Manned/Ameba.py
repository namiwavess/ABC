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
    self.mousex:mouse's coordinate x(ウィンドウの中での)
    self.mousey:mouse's coordinate y(ウィンドウの中での)
    """
    super().__init__(name, url)

  def initXY(self):
    self.x = 5#random.randint(0, self.width)
    self.y = 5#random.randint(0, self.height)

  def cal_mousedistance(self):
    return  int(math.sqrt(math.pow(self.mousex - self.x, 2) + math.pow(self.mousey - self.y, 2)))

  def cal_vector(self):
    angle = math.atan2(self.mousey - self.y, self.mousex - self.x)

    if angle < 0:
      angle += 2 * math.pi

    speed = self.cal_mousedistance()
    if speed > 3:
      speed = 3
    vx = speed * math.cos(angle)
    vy = speed * math.sin(angle)
    
    return int(vx), int(vy)


  def play(self):
    '''
    #--------dbug-------------------
    print("\u001B[0K",end="")      
    print(" ")
    print("\u001B[0K",end="")
    print("---------------------------------")
    print("\u001B[0K",end="")
    print("mouse[x,y]:[",self.mousex,",",self.mousey,"]")
    print("\u001B[0K",end="")
    print("mass :",self.mass)
    print("\u001B[0K",end="")
    print(" now :[",self.x,",",self.y,"]")
    print("\u001B[0K",end="")
    print("---------------------------------")
    print("\u001B[6A",end="")
    #-------------------------------------------
    '''
    return self.cal_vector()

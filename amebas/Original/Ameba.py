from .Cell import Cell
import random
import math
import numpy as np
import time
import copy
import sys

class Ameba(Cell):
  def __init__(self, name, url):
    """
    self.x: ameba's position x
    self.y: ameba's position y
    self.mass: ameba's radius
    self.width = 0: world width
    self.height = 0: world height
    self.food: food positions
    self.players:[(number, x, y, mass)] (include yourself)
    """
    self.printcnt = 0
    self.calc = True
    self.routeNum = 3
    self.route = [(-1,-1),(-1,-1),(-1,-1),-1]
    self.preX = 0
    self.preY = 0
    super().__init__(name, url)

  def cal_distance(self, x, y):
    return math.sqrt(math.pow(x - self.x, 2) + math.pow(y - self.y, 2))

  def cal_two_distance(self, x1, y1, x2, y2):
    return math.sqrt(math.pow(x1 - x2, 2) + math.pow(y1 - y2, 2))

  def cal_vector(self, dX, dY):
    angle = math.atan2(dY - self.y, dX - self.x)

    if angle < 0:
      angle += 2 * math.pi

    speed = 3
    vx = speed * math.cos(angle)
    vy = speed * math.sin(angle)

    return int(vx), int(vy)

  def select_candidate(self,pos,foodlist):
    so = sorted(foodlist,key=lambda x:self.cal_two_distance(x[0],x[1],pos[0],pos[1]))
    return so[0],so[1],so[2]

  def find_max_density(self):
    quantity = np.zeros((11,11))
    w = int(self.width/10)
    h = int(self.height/10)

    for fo in self.food:
      hor = int(fo[0]/w)
      ver = int(fo[1]/h)
      quantity[ver][hor] += 1

    ma =[0,0]

    for row in range(11):
      for col in range(11):
        if quantity[ma[0]][ma[1]] < quantity[row][col]:
          ma[0] = row
          ma[1] = col

    retx = int(w*ma[0] + w/2)
    rety = int(h*ma[1] + h/2)

    #print("w:",w,",h:",h,",hor:",hor,",ver:",ver,",max:",ma,",retx:",retx,",rety:",rety)
    #print(quantity)

    return retx,rety


  def initXY(self):
    #self.x = random.randint(0, self.width)
    #self.y = random.randint(0, self.height)
    self.x,self.y=self.find_max_density()
    #print("init :",self.x,",",self.y)

  def cal(self):

    pos = [self.x,self.y]
    temp1 = self.select_candidate(pos,self.food)
    candidate =[]
    distance = [0 for i in range(9)]
    index = 0
    for i in range(3):

      pos1 = temp1[i]
      foodlist1 = copy.deepcopy(self.food)
      foodlist1.remove(temp1[i])

      dist_temp = self.cal_two_distance(pos[0],pos[1],pos1[0],pos1[1])
      temp2 = self.select_candidate(pos1,foodlist1)
      for j in range(3):

        pos2 = temp2[j]
        foodlist2 = copy.deepcopy(foodlist1)
        foodlist2.remove(temp2[j])

        distance[index] += self.cal_two_distance(pos1[0],pos1[1],pos2[0],pos2[1])
        
        temp3 = sorted(foodlist2,key=lambda x:self.cal_two_distance(pos2[0],pos2[1],x[0],x[1]))
        
        candidate.append([temp1[i],temp2[j],temp3[0]])
        distance[index] += (dist_temp+self.cal_two_distance(pos2[0],pos2[1],temp3[0][0],temp3[0][1]))
        index += 1

    p=[]
    k = 0
    for i,j in zip(candidate,distance):
      i.append(j)
      p.append(i)
      k+=1

    p = sorted(p,key=lambda x:x[3])

    return p[0]


  def play(self):
    #if self.printcnt==100:
    #  self.printcnt=0
    #self.printcnt += 1
  
    if self.route[self.routeNum] not in self.food:
      self.routeNum += 1

    if self.routeNum >= 3:
      self.routeNum = 0
      self.calc = True

    if self.calc:
      self.route = self.cal()
      self.routeNum
      self.calc = False

    distination = self.route[self.routeNum]
    direct = self.cal_vector(distination[0], distination[1])
    '''
    #--------dbug-------------------
    if self.printcnt > -1:
      print("\u001B[0K",end="")      
      print(" ")
      print("\u001B[0K",end="")
      print("---------------------------------")
      print("\u001B[0K",end="")
      print("route:",self.route)
      print("\u001B[0K",end="")
      print(" num :",self.routeNum)
      print("\u001B[0K",end="")
      print("dist :",distination)
      print("\u001B[0K",end="")
      print("direc:",direct)
      print("\u001B[0K",end="")
      print("mass :",self.mass)
      print("\u001B[0K",end="")
      print(" now :[",self.x,",",self.y,"]")
      print("\u001B[0K",end="")
      print("diff :[",self.x-self.preX,",",self.y-self.preY,"]")
      self.preX = self.x
      self.preY = self.y
      print("\u001B[0K",end="")
      print("---------------------------------")
      print("\u001B[11A",end="")

      self.printcnt = 0
    else:
      self.printcnt += 1
    #-------------------------------------------
    '''

    return direct
    #return 0,0
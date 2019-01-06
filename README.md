# ABC
Cell prey game inspired from Agar.io

## Todo
 - amebaAI
 - UI
 - 対人戦
 - Debug
 - パラメータ調整

## Feild
- 1 frame40 milli-sec

## Ameba
- ameba同士の衝突はmassが5%以上大きい方に吸収
- amebaのmass（半径）は最小2，最大が500
- 餌か他のamebaを吸収しないと，フレーム毎にmassが0.2%減
- 速度はmax 3。 dx^2 + dy^2 <= 9

## Protocol
- 1 userid food width height
    - 接続してきたプレイヤーに固有のuseridとそのworldのfoodの初期位置，worldの横と縦
- 2 userid name x y
    - プレイヤーは自分のuseridと名前，初期位置(x,y)をサーバに伝える
- 3 userid dx dy
    - プレイヤーは自分の行きたい方向(dx, dy)をuseridとともに伝える
- 4 diff amebas
    - 1 frame毎にプレイヤー全体に前回のfoodのdiffとすべてのプレーヤの位置が送られる

## Install
pip3 install websocket-client==0.48 --user

## Run
./main -num 10
python3 run.py 10

## Compile
env GOOS=linux GOARCH=amd64 go build main.go

## Docker
    docker build ./ -t ameba
    
## Demo
準備中

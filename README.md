# AmebaBattleCircle
Agar.io (local version)

## 課題
 - 強いamebaの作成
 - ABCの改良
    - 表示の改良
    - 人間が戦えるようにする
    - デバッグ
    - パラメータ調整

## Feildの制限

- 1 frameは40 milli-sec

## Amebaの設定

- ameba同士がぶつかった場合，massが5%以上大きいほうに吸収される．
- amebaのmass（半径）は最小2，最大が500
- 餌か他のamebaを吸収しないと，フレーム毎にmassが0.2%減る．
- 速度はmax 3．つまり dx^2 + dy^2 <= 9

## プロトコル

- 1 userid food width height
    - 接続してきたプレイヤーに固有のuseridとそのworldのfoodの初期位置，worldの横と縦
- 2 userid name x y
    - プレイヤーは自分のuseridと名前，初期位置(x,y)をサーバに伝える
- 3 userid dx dy
    - プレイヤーは自分の行きたい方向(dx, dy)をuseridとともに伝える
- 4 diff amebas
    - 1 frame毎にプレイヤー全体に前回のfoodのdiffとすべてのプレーヤの位置が送られる


## 準備

pip3 install websocket-client==0.48 --user


## 起動

./main -num 10
python3 run.py 10


## コンパイル

env GOOS=linux GOARCH=amd64 go build main.go

# Dockerで起動

    docker build ./ -t ameba
    docker run --rm -it -p 3000:3000 ameba

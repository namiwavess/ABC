# ABC
Cell prey game inspired from Agar.io

## Todo
 - amebaAI
 - UI
 - 対人戦
 - Debug
 - パラメータ調整

## Feildの制限

- 1 frame 40 ms

## Ameba

- ameba同士が衝突した場合、mass(半径)が5%以上大きい方に吸収される。
- amebaのmassは 最小2 ~ 最大500
- 餌か他のamebaを吸収しないと、フレーム毎にmassが0.2%減る。
- 速度は max 3 (dx^2 + dy^2 <= 9)

## Protocol

- 1 userid food width height
    - 接続してきたプレイヤーに固有のuseridとそのworldのfoodの初期位置，worldの横と縦
- 2 userid name x y
    - プレイヤーは自分のuseridと名前，初期位置(x,y)をサーバに伝える
- 3 userid dx dy
    - プレイヤーは自分の行きたい方向(dx, dy)をuseridと共に伝える
- 4 diff amebas **mousex mousey**
    - 1 frame毎にプレイヤー全体に前回のfoodのdiffとすべてのプレーヤの位置が送られる
- 5　mousex mousey
    - ブラウザからマウスの位置をサーバに送信
- 6 message
    - gamelogの表示するメッセージを送る
- 7　"Quit,game"
    - サーバをキルしてゲーム終了
- 8　message
    - amebaに吸収された場合、または、優勝したときにメッセージ(キル数、順位など)を送信
- 9 未実装

サーバを立てる際に``-c True``を付けると自分のsizeが常に200以下にはならない仕様に。

## Install

pip3 install websocket-client==0.48 --user

## Run

./main -num 10
python3 run.py 10

## Compile

env GOOS=linux GOARCH=amd64 go build main.go

# Docker

docker build ./ -t ameba
docker run --rm -it -p 3000:3000 ameba

## Demo

https://www.youtube.com/watch?v=WSIeLhvBBtA

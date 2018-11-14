import * as PIXI from 'pixi.js'
import { AnySoaRecord } from 'dns';
import { truncate } from 'fs';

var url = "ws://" + window.location.host + "/view"
var url = "ws://localhost:3000/view"
var ws = new WebSocket(url)

export default class AmebaWorld extends PIXI.Application
{
	food:PIXI.Graphics[] = []
	players:PIXI.Graphics[] = []

	constructor()
	{
		super({
			view: <HTMLCanvasElement>document.getElementById('canvas'),
			backgroundColor: 0xf8f8f8,
			width: 1000,
			height: 1000
		});

		document.body.appendChild(this.view);

		for (var i = 0; i < 100; i++) {
			this.food.push(new PIXI.Graphics())
			this.food[i].beginFill(0xe49758)
			this.food[i].drawEllipse(-1, -1, 2, 2)
			this.food[i].endFill()
			this.food[i].visible = false
			this.stage.addChild(this.food[i])
		}

		for (var i = 0; i < 100; i++) {
			this.players.push(new PIXI.Graphics())
			this.players[i].beginFill(0xff0000)
			this.players[i].drawEllipse(-1, -1, 10, 10)
			this.players[i].endFill()
			this.players[i].visible = false
			this.stage.addChild(this.players[i])
		}

		for (var i = 0; i < 1000; i += 50) {
			var line = new PIXI.Graphics()
			line.lineStyle(1, 0xc6c6c6).moveTo(0, i).lineTo(1000, i)
			this.stage.addChild(line)
			line = new PIXI.Graphics()
			line.lineStyle(1, 0xc6c6c6).moveTo(i, 0).lineTo(i, 1000)
			this.stage.addChild(line)
		}


		this.ticker.add((deltaTime) => this.update(deltaTime));
	}

	updateFood(pos:number[]) {
		let i, idx:number

		for (i = 0; i < pos.length; i += 3) {
			idx = pos[i]
			this.food[idx].x = pos[i + 1]
			this.food[idx].y = pos[i + 2]
			if (this.food[idx].x != -1) {
				this.food[idx].visible = true
			}
			else {
				this.food[idx].visible = false
			}
		}
	}

	updatePlaysers(pos:number[]) {
		let i, j:number
		for (i = 0, j = 0 ; i < this.players.length; i++) {
			if (i == pos[j]) {
				this.players[i].x = pos[j + 1]
				this.players[i].y = pos[j + 2]
				this.players[i].width = pos[j + 3] * 2
				this.players[i].height = pos[j + 3] * 2
				this.players[i].visible = true
				j += 4
			}
			else {
				this.players[i].visible = false
			}

		}
	}

	update(deltaTime:number)
	{
	}
}

let world = new AmebaWorld()

ws.onmessage = (msg) => {
	let cmd = msg.data.split(" ")
	switch (cmd[0]) {
		case "4":
		let f = cmd[1].split(",").map((elem) => {return parseInt(elem)}).slice(0, -1)
		let p = cmd[2].split(",").map((elem) => {return Number(elem)}).slice(0, -1)
		world.updatePlaysers(p)
		world.updateFood(f)
		break
		default:
	}
}

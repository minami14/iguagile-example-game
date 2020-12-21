const width = 800;
const height = 600;
const app = new PIXI.Application({
    width: width, height: height, backgroundColor: 0x1099bb, resolution: window.devicePixelRatio || 1,
});

document.body.appendChild(app.view);

const cardBack = PIXI.Texture.from('img/card_back.png');
const cardFront = [];
for (let i = 0; i < 10; i++) {
    cardFront[i] = PIXI.Texture.from(`img/card${i}.png`);
}

function onClick() {
    console.log(this);
    if (current === null || player === null) {
        return;
    }
    if (current !== player.id) {
        return;
    }
    const req = {
        card_number: this.number,
    }
    const data = JSON.stringify(req);
    Game.client.send(data);
}

cards = [];
for (let i = 0; i < 20; i++) {
    const card = new PIXI.Sprite(cardBack);
    card.interactive = true;
    card.anchor.set(0.5);
    card.on('pointerdown', onClick);
    card.x = ((i % 5) + 1) * width / 6;
    card.y = (Math.floor(i / 5) + 1) * height / 5
    card.number = i;
    app.stage.addChild(card);
    cards.push(card);
}

const style = new PIXI.TextStyle({
    fill: "white",
    strokeThickness: 2,
    fontSize: 32,
});
const text = new PIXI.Text('Waiting for an opponent to connect', style);
app.stage.addChild(text);
let player;
let current;

Game.start(function () {
    Game.client.onConnect = function () {
        const p = document.createElement('p');
        const room = {
            room_id: Game.room.room_id,
            server: Game.room.server,
            application_name: Game.room.application_name,
            version: Game.room.version,
            password: Game.room.password,
        }
        const joinToken = btoa(JSON.stringify(room));
        const joinUrl = `${document.location.origin}/game/join.html?token=${joinToken}`;
        const a = document.createElement('a');
        a.href = joinUrl;
        a.target = '_blank'
        a.appendChild(document.createTextNode(joinUrl));
        p.appendChild(document.createTextNode(`参加者用URL: `));
        p.appendChild(a);
        document.body.append(p);
    };
    Game.client.onReceive = function (data) {
        console.log(data);
        data = JSON.parse(data);
        let card;
        switch (data.message_type) {
            case "start":
                player = data.player;
                break;
            case "open":
                card = cards[data.card.number];
                card.texture = cardFront[data.card.id];
                break;
            case "close":
                cards.forEach((c, i) => {
                    c.texture = cardBack;
                });
                break;
            case "get":
                data.cards.forEach((c, i) => {
                    app.stage.removeChild(cards[c.number]);
                });
                break;
            case "change":
                current = data.player.id;
                if (current === player.id) {
                    text.text = 'Your turn'
                } else {
                    text.text = 'Opponent\'s turn'
                }
                break;
            case "finish":
                if (data.player.id === player.id) {
                    text.text = 'Win!'
                } else {
                    text.text = 'Lose…'
                }
                break;
            case "error":
                console.log(data);
                break;
            default:
                console.log(data);
                break;
        }
    };
    Game.client.connect(Game.room);
});

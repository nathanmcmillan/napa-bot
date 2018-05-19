
const lines = { canvas: null, context: null, timeSlices: 30, hourInterval: 1 };
let canvas = document.createElement('canvas');
canvas.style.display = 'block';
canvas.style.position = 'absolute';
canvas.style.left = '0';
canvas.style.right = '0';
canvas.style.top = '0';
canvas.style.bottom = '0';
canvas.style.margin = 'auto';
canvas.width = window.innerWidth * 0.5;
canvas.height = window.innerHeight * 0.25;
let context = canvas.getContext('2d');
lines.canvas = canvas;
lines.context = context;
drawLineGraph();
document.body.appendChild(canvas);

const status = document.getElementById('status');
const out = document.getElementById('out');
const socket = new WebSocket('ws://localhost:80/websocket');

socket.onopen = function (event) {
    status.innerHTML = 'listening';
};

socket.onmessage = function (event) {
    if (event.data === null) {
        return;
    }
    let js = JSON.parse(event.data);
    switch (js.uid) {
        case 'ticker':
            out.innerHTML += '<p>' + js.uid + ', ' + js.time + ', ' + js.product_id + ', ' + js.price + ', ' + js.side + '</p>';
            break;
        case 'log':
            out.innerHTML += '<p>' + js.message + '</p>';
            break;
    }
};

socket.onerror = function (event) {
    status.innerHTML += '<p>error ' + event.data + '</p>';
};

socket.onclose = function () {
    status.innerHTML = 'closed';
    socket = null;
};

function buy() {
    if (no()) return;
    socket.send(`{"uid":"buy", "BTC-USD"}`);
}

function subscribe() {
    if (no()) return;
    socket.send(`{"uid":"sub-exchange"}`);
}

function unsubscribe() {
    if (no()) return;
    socket.send(`{"uid":"unsub-exchange"}`);
}

function no() {
    if (socket === null) {
        status.innerHTML = 'must be connected to napa server';
        return true;
    }
    return false;
}

function drawLineGraph() {
    let canvas = lines.canvas;
    let context = lines.context;
    context.fillStyle = 'rgb(100, 100, 100)';
    context.fillRect(0, 0, canvas.width, canvas.height);

    context.lineWidth = 2;
    context.strokeStyle = 'rgb(200, 100, 0)';
    context.beginPath();
    context.moveTo(10, 10);
    context.lineTo(100, 100);
    context.stroke();

    let begin = Date.now();
    begin.setDate(begin.getDate() - 7); // one week
    let end = Date.now();
}
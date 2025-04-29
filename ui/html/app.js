const ws = new WebSocket(`ws://${window.location.host}/ws`);

const messagesDiv = document.getElementById('messages');
const messageInput = document.getElementById('messageInput');
const sendButton = document.getElementById('sendButton');
const switchButton = document.getElementById('switchButton');

let currentProtocol = "tcp";

ws.onopen = () => {
    addMessage("Connected to server using TCP");
};

ws.onmessage = (event) => {
    addMessage(`Server: ${event.data}`);
};

ws.onclose = () => {
    addMessage("Disconnected from server");
};

sendButton.onclick = () => {
    const msg = messageInput.value;
    if (msg && ws.readyState === WebSocket.OPEN) {
        ws.send(msg);
        addMessage(`You: ${msg}`);
        messageInput.value = "";
    }
};

switchButton.onclick = () => {
    currentProtocol = currentProtocol === "tcp" ? "udp" : "tcp";
    ws.send(`switch:${currentProtocol}`);
    switchButton.textContent = currentProtocol === "tcp" ? "Switch to UDP" : "Switch to TCP";
    addMessage(`Switched to ${currentProtocol.toUpperCase()}`);
};

function addMessage(msg) {
    const p = document.createElement('p');
    p.textContent = msg;
    messagesDiv.appendChild(p);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}

let latencySpan = document.getElementById("latency");

function pingServer() {
    const timestamp = Date.now();
    ws.send(`__ping__:${timestamp}`);
}

setInterval(pingServer, 3000); // Ping every 3 seconds

ws.onmessage = (event) => {
    const data = event.data;

    if (data.startsWith("__pong__:")) {
        const sentTime = parseInt(data.split(":")[1], 10);
        const rtt = Date.now() - sentTime;
        latencySpan.textContent = rtt;
        return;
    }

    addMessage(`Server: ${data}`);
};

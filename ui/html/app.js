let protocol = "tcp";
const ws = new WebSocket("ws://localhost:8080/ws");

document.getElementById("sendBtn").onclick = () => {
    const msg = document.getElementById("message").value;
    ws.send(msg);
    document.getElementById("message").value = "";
};

document.getElementById("switchBtn").onclick = () => {
    protocol = (protocol === "tcp") ? "udp" : "tcp";
    ws.send("switch:" + protocol);
    document.getElementById("switchBtn").innerText = "Switch to " + (protocol === "tcp" ? "UDP" : "TCP");
};

ws.onmessage = (event) => {
    const chatbox = document.getElementById("chatbox");
    chatbox.innerHTML += `<div>${event.data}</div>`;
};

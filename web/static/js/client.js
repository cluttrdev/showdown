var ws = new WebSocket("ws://" + window.location.host + "/ws");

ws.onmessage = (event) => {
    var msg = JSON.parse(event.data);
    switch (msg.type) {
        case 'title':
            document.querySelector('title').innerHTML = msg.data;
            break;
        case 'content':
            document.querySelector('.markdown-body').innerHTML = msg.data;
            break;
    }
}

ws.onerror = (event) => {
    console.debug(event)
}


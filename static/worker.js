var ws = new WebSocket("ws://162.245.217.17:8080/game");

ws.onopen = function() {
  console.log("connected\n");
}
ws.onmessage = function(evt){
  postMessage(evt.data);
  console.log("[ALERT] Data Recieved: " + evt.data + "\n");
}

// Log errors
ws.onerror = function (error) {
  console.log('[Error] websocket Error:  ' + error + "\n");
};

onmessage = function(e){
  console.log("sending the following command: " + e.data[0] + "\n");
  switch(e.data[0])
  {
    case 37: //left
    ws.send('{"command":"MOVE_LEFT", "id":"' + e.data[1] + '"}');
    break;

    case 38: //up
    ws.send('{"command":"MOVE_UP", "id":"' + e.data[1] + '"}');
    break;

    case 39: //right
    ws.send('{"command":"MOVE_RIGHT", "id":"' + e.data[1] + '"}');
    break;

    case 40: //down
    ws.send('{"command":"MOVE_DOWN", "id":"' + e.data[1] + '"}');
    break;
  }
}

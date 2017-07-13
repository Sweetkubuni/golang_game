var ws = new WebSocket("ws://162.245.217.17:8080/game");
ws.onmessage = function(evt){
  postMessage(evt.data);
}

onmessage = function(e){
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

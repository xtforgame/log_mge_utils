(function (window) {
  var lmu = window.lmu = {};
  lmu.init = function(onlyLogger) {
    var EventNextIterationCode = String.fromCharCode({{.EventNextIterationCode}});
    var EventOnDataCode = String.fromCharCode({{.EventOnDataCode}});
    var EventLogRemovedCode = String.fromCharCode({{.EventLogRemovedCode}});

    lmu.loggerField = document.getElementById("logger-field");
    lmu.loggerField.innerHTML = `
      <fieldset>
        <legend>Logger</legend>
        <button id="logger-open">Open</button>
        <button id="logger-close">Close</button>
        <br />
        <input id="logger-input" type="text" value="Hello world!" />
        <button id="logger-send">Send</button>
        <button id="logger-clear">Clear</button>
        <button id="logger-removelog">Remove Log</button>
      </fieldset>
    `;
    lmu.output = document.getElementById("ws-output");
    lmu.input = document.getElementById("logger-input");

    lmu.currentSendElem = null;
    lmu.currentResponseElem = null;

    document.getElementById("logger-open").onclick = function (evt) {
      if (lmu.readerWs) {
        return false;
      }
      if (!onlyLogger) {
        lmu.readerWs = new WebSocket("{{.WsBaseUrl}}/listeners/20022");
        lmu.writerWs = new WebSocket("{{.WsBaseUrl}}/loggers/20022");
      } else {
        lmu.readerWs = new WebSocket("{{.WsBaseUrl}}/loggers/20022");
        lmu.writerWs = lmu.readerWs;
      }
      lmu.readerWs.binaryType = 'blob'
      lmu.readerWs.onopen = function (evt) {
        lmu.open();
        lmu.readerWs.send("password");
      }
      lmu.readerWs.onclose = function (evt) {
        lmu.print("READER CLOSE");
        if (lmu.readerWs === lmu.writerWs) {
          lmu.writerWs = null;
        }
        lmu.readerWs = null;
      }
      lmu.readerWs.onmessage = function (evt) {
        if (evt.data instanceof Blob) {
          evt.reader = new FileReader();
          evt.reader.onload = function () {
            if (evt.reader.result && evt.reader.result.length > 0) {
              console.log(evt.reader.result);
              switch (evt.reader.result[0]) {
                case EventOnDataCode:
                  lmu.appendResponse(evt.reader.result.slice(2));
                  break;
                case EventNextIterationCode:
                  lmu.clearResponse();
                  break;
                case EventLogRemovedCode:
                  lmu.logRemoved();
                  break;
              }
            }
          };
          evt.reader.readAsText(evt.data);
        } else {
          lmu.print("RESPONSE: " + evt.data);
        }
      }
      lmu.readerWs.onerror = function (evt) {
        lmu.print("READER ERROR: " + evt.data);
      }
      if (lmu.writerWs !== lmu.readerWs) {
        lmu.writerWs.binaryType = 'blob'
        lmu.writerWs.onopen = function (evt) {
          // lmu.open();
          lmu.writerWs.send("password");
        }
        lmu.writerWs.onclose = function (evt) {
          lmu.print("WRITER CLOSE");
          lmu.writerWs = null;
        }
        lmu.writerWs.onerror = function (evt) {
          lmu.print("WRITER ERROR: " + evt.data);
        }
      }
      return false;
    };
    document.getElementById("logger-send").onclick = function (evt) {
      if (!lmu.writerWs) {
        return false;
      }
      lmu.displaySend("SEND: " + lmu.input.value);
      lmu.writerWs.send("\1" + lmu.input.value);
      return false;
    };
    document.getElementById("logger-clear").onclick = function (evt) {
      if (!lmu.writerWs) {
        return false;
      }
      lmu.writerWs.send("\2");
      return false;
    };
    document.getElementById("logger-removelog").onclick = function (evt) {
      if (!lmu.writerWs) {
        return false;
      }
      lmu.writerWs.send("\3");
      return false;
    };
    document.getElementById("logger-close").onclick = function (evt) {
      if (!lmu.readerWs) {
        return false;
      }
      lmu.readerWs.close();
      if (lmu.writerWs !== lmu.readerWs) {
        lmu.writerWs.close();
      }
      return false;
    };
  };

  lmu.newDiv = function (innerHTML) {
    var d = document.createElement("div");
    if (innerHTML) {
      d.innerHTML = innerHTML;
    }
    return d;
  };

  lmu.open = function () {
    lmu.output.innerHTML = '';
    lmu.output.appendChild(lmu.newDiv('OPEN'));

    lmu.currentSendElem = lmu.newDiv();
    lmu.output.appendChild(lmu.currentSendElem);

    lmu.currentResponseElem = lmu.newDiv('RESPONSE: ');
    lmu.currentResponseElem.setAttribute('style', 'word-wrap: break-word; white-space: pre-wrap; word-break: break-all; font-family: none;');
    lmu.output.appendChild(lmu.currentResponseElem);
  };

  lmu.print = function (message) {
    lmu.output.appendChild(lmu.newDiv(message));
    // lmu.currentResponseElem = null;
    // lmu.currentSendElem = null;
  };

  lmu.displaySend = function (message) {
    if (!lmu.currentSendElem) {
      return;
    }
    lmu.currentSendElem.innerHTML = message;
  };

  lmu.appendResponse = function (message) {
    if (!lmu.currentResponseElem) {
      return;
    }
    lmu.currentResponseElem.innerHTML += message;
  };

  lmu.logRemoved = function () {
    if (!lmu.currentResponseElem) {
      lmu.currentResponseElem = document.createElement("pre");
      lmu.currentResponseElem.setAttribute("style", "word-wrap: break-word; white-space: pre-wrap; word-break: break-all; font-family: none;");
      lmu.currentResponseElem.innerHTML = "RESPONSE: ";
      lmu.output.appendChild(lmu.currentResponseElem);
    }
    lmu.currentResponseElem.innerHTML += "<br />===log removed===<br />";
  };

  lmu.clearResponse = function () {
    lmu.displaySend("")
    lmu.appendResponse("")
    lmu.currentResponseElem.innerHTML = "RESPONSE: ===new iteration===<br />";
  };
}(window));

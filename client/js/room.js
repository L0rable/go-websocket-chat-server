window.onload = function () {
    var conn;
    const msg = document.getElementById("msg");
    const msgLog = document.getElementById("msg-log");

    const urlParams = new URLSearchParams(window.location.search);
    const clientId = urlParams.get("clientId");
    const clientName = urlParams.get("clientName");
    const roomNo = urlParams.get("roomNo");

    if (roomNo != null) {
        document.title = "Room " + roomNo;
    } else {
        document.title = "Room ";
    }

    document.getElementById("msg-send-form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }
        conn.send(msg.value);
        msg.value = "";
        return false;
    };

    document.getElementById("leave-input").addEventListener("click", async function (e) {
        e.preventDefault();

        if (!clientName || !roomNo) {
            return false;
        }

        const leaveReqData = JSON.stringify({
            clientId: clientId,
            clientName: clientName,
            room: Number(roomNo),
        });

        try {
            const leaveReq = await fetch('/leave', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: leaveReqData
            });

            if (!leaveReq.ok) {
                throw new Error("Failed to leave chat room");
            }

            if (leaveReq.redirected) {
                window.location.href = leaveReq.url;
            }
        } catch (err) {
            console.error("Error:", err);
        }
    });

    if (window["WebSocket"]) {
        conn = new WebSocket(`wss://${window.location.host}/ws?clientId=${clientId}&clientName=${clientName}&roomNo=${roomNo}`);
        conn.onopen = function (evt) {
            console.log("ws conn open");
            if (evt.data == 'null') {
                return;
            }
        };

        conn.onclose = function (evt) {
            console.log("ws conn close");
            WebSocket.close();
            var msgText = document.createElement("div");
            msgText.innerHTML = "<b>Connection closed.</b>";
            msgLog.appendChild(msgText)
        };

        conn.onmessage = function (evt) {
            if (evt.data == 'null') {
                return;
            }

            var msgs = evt.data;
            var msgsJSON = JSON.parse(msgs);
            if (!Array.isArray(msgsJSON)) {
                msgsJSON = [msgsJSON];
            }
            msgsJSON.forEach(msgJSON => {
                const clientName = msgJSON.clientName;
                const text = msgJSON.text;
                var msgText = document.createElement("div");
                msgText.innerText = `${clientName}: ${text}`;
                msgLog.appendChild(msgText);
            });
        };

    } else {
        var msgText = document.createElement("div");
        msgText.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        msgLog.appendChild(item);
    }
};
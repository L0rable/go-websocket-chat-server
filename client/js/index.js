window.onload = function () {
    const name = document.getElementById("user-name");
    const room = document.getElementById("chat-room-no");

    const urlParams = new URLSearchParams(window.location.search);
    const clientId = urlParams.get("clientId");
    const clientName = urlParams.get("clientName");

    if (clientId != null) {
        document.getElementById("user-info-clientId").innerHTML = "UUID: " + clientId;
    }
    if (clientName != null) {
        document.getElementById("user-info-clientName").innerHTML = "Name: " + clientName;
    }

    document.getElementById("room-join-form").onsubmit = async function (e) {
        e.preventDefault();
        
        if (!name.value || !room.value) {
            return false;
        }
        const joinReqData = JSON.stringify({
            clientId: clientId,
            clientName: name.value,
            room: Number(room.value),
        });

        try {
            const joinReq = await fetch('/join', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: joinReqData
            });

            if (!joinReq.ok) {
                throw new Error("Failed to join chat room");
            }

            if (joinReq.redirected) {
                window.location.href = joinReq.url;
            }
        } catch (err) {
            console.error("Error:", err);
        }

        return false;
    };
};

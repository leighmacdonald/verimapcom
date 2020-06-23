import {EvtMessage, EvtSetMission, WSClient} from "./ws";
import {map_view_mission} from "./maps";

class payloadMessage {
    constructor() {
        this.mission_id = 0;
        this.person_name = "";
        this.person_id = 0;
        this.message = "";
    }
}


function send_message(ws, message) {
    console.log("Sending", message)
    let p = new payloadMessage()
    p.message = message
    ws.send(EvtMessage, p)
}

function recv_message(payload) {
    console.log("got msg:", payload)
}

class payloadSetMission {
    constructor(mission_id) {
        this.mission_id = mission_id;
    }
}
/**
 *
 * @param ws WSClient
 * @param mission_id
 */
function set_mission(ws, mission_id) {
    mission_id = parseInt(mission_id, 10)
    if (mission_id > 0) {
        ws.send(EvtSetMission, new payloadSetMission(mission_id))
    }
}

export function page_mission(mission_id) {
    let ws = new WSClient();
    ws.onopen = (event) => {
        console.log("Calling set_mission")
        set_mission(ws, mission_id)
    }
    const submit = document.getElementById("chat_submit")
    const msg = document.getElementById("chat_message")
    submit.addEventListener("click", (e) => {
        e.preventDefault()
        let user_message = msg.value
        send_message(ws, user_message);
        msg.value = "";
    })
    ws.register(EvtMessage, recv_message)
    const map_div = document.getElementById("map");
    map_view_mission("map", mission_id, map_div.dataset.lat_ul, map_div.dataset.lon_ul)
}
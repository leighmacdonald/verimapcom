import {map_view_mission} from "./maps";
import {RPCClient} from "./pb/RpcServiceClientPb";
import {open_mission} from "./client";

function send_message(message) {
    const element = <HTMLDivElement>document.getElementById("live_chat")
    add_event_log(element, {"payload": {"message": message}})
}

function recv_message(payload) {
    console.log("got msg:", payload)
}

/**
 *
 * @param ws WSClient
 * @param mission_id
 */
function set_mission(mission_id: number) {
    if (mission_id > 0) {

    }
}

function add_event_log(element: HTMLDivElement, event) {
    const d = `
        <div class="cell">
            <div class="grid-x">
                <div class="cell small-2">
                    <a href="/profile/${event.payload["profile_id"]}">${event.payload["profile_id"]}</a>
                </div>
                <div class="cell auto">
                    <p>${event.payload["message"]}</p>
                </div>
            </div>
        </div>
    `

    element.insertAdjacentHTML('beforeend', d);
}

function fetch_mission_events(mission_id: number) {
    const element = <HTMLDivElement>document.getElementById("live_chat")
    fetch(`/mission/${mission_id}/events`)
        .then(response => response.json())
        .then(data => {
            data.forEach((event) => {
                add_event_log(element, event)
            })
        })
}

export function page_mission(mission_id: number) {
    const svc = new RPCClient('http://172.16.1.4:8800', null, null);
    open_mission(svc, mission_id);
    fetch_mission_events(mission_id)

    const submit = document.getElementById("chat_submit")
    const msg = <HTMLInputElement>document.getElementById("chat_message")
    submit.addEventListener("click", (e) => {
        e.preventDefault()
        let user_message = msg.value;
        //send_message(ws, user_message);
        msg.value = "";
    })
    const map_div = document.getElementById("map");
    map_view_mission("map", mission_id, map_div.dataset.lat_ul, map_div.dataset.lon_ul)
}
import {map_view_mission} from "./maps";
import {Client, wsEvent} from "./client";

function send_message(message) {
    const element = <HTMLDivElement>document.getElementById("live_chat")
    add_event_log(element, {"payload": {"message": message}})
}

function recv_message(payload) {
    console.log("got msg:", payload)
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

function handle_open_mission(payload: any) {
    console.log("got open mission payload", payload)
}

function handle_mission_events(payload: any) {

}

function handle_new_flight(payload: any) {

}

function handle_position(payload: any) {

}

function handle_message_recv(payload: any) {

}

function ui_send_message(c: Client, e: Event) {
    const msg = <HTMLInputElement>document.getElementById("chat_message")
    e.preventDefault()
    let user_message = msg.value;
    c.message_send(user_message)
    msg.value = "";
}

export function page_mission(mission_id: number) {
    const client = new Client();
    client.register(wsEvent.missionOpen, handle_open_mission)
    client.register(wsEvent.missionEvents, handle_mission_events)
    client.register(wsEvent.missionNewFlight, handle_new_flight)
    client.register(wsEvent.missionPosition, handle_position)
    client.register(wsEvent.missionRecvMessage, handle_message_recv)

    const submit = document.getElementById("chat_submit")
    submit.addEventListener("click", (e) => {
        ui_send_message(client, e)
    })

    const map_div = document.getElementById("map");
    map_view_mission("map", mission_id, map_div.dataset.lat_ul, map_div.dataset.lon_ul)

    client.open_mission(mission_id)
}

"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.page_mission = void 0;
var maps_1 = require("./maps");
var rpc_grpc_web_pb_1 = require("./pb/rpc_grpc_web_pb");
function send_message(ws, message) {
    var element = document.getElementById("live_chat");
    add_event_log(element, { "payload": { "message": message } });
}
function recv_message(payload) {
    console.log("got msg:", payload);
}
/**
 *
 * @param ws WSClient
 * @param mission_id
 */
function set_mission(ws, mission_id) {
    mission_id = parseInt(mission_id, 10);
    if (mission_id > 0) {
    }
}
function add_event_log(element, event) {
    var d = "\n        <div class=\"cell\">\n            <div class=\"grid-x\">\n                <div class=\"cell small-2\">\n                    <a href=\"/profile/" + event.payload["profile_id"] + "\">" + event.payload["profile_id"] + "</a>\n                </div>\n                <div class=\"cell auto\">\n                    <p>" + event.payload["message"] + "</p>\n                </div>\n            </div>\n        </div>\n    ";
    element.insertAdjacentHTML('beforeend', d);
}
function fetch_mission_events(mission_id) {
    var element = document.getElementById("live_chat");
    fetch("/mission/" + mission_id + "/events")
        .then(function (response) { return response.json(); })
        .then(function (data) {
        data.forEach(function (event) {
            add_event_log(element, event);
        });
    });
}
function page_mission(mission_id) {
    var svc = new rpc_grpc_web_pb_1.RPCClient('http://172.16.1.4:8800', null, null);
    fetch_mission_events(mission_id);
    var submit = document.getElementById("chat_submit");
    var msg = document.getElementById("chat_message");
    submit.addEventListener("click", function (e) {
        e.preventDefault();
        var user_message = msg.value;
        //send_message(ws, user_message);
        msg.value = "";
    });
    var map_div = document.getElementById("map");
    maps_1.map_view_mission("map", mission_id, map_div.dataset.lat_ul, map_div.dataset.lon_ul);
}
exports.page_mission = page_mission;

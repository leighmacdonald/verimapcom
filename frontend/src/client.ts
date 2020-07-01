import {has, forEach, map} from "lodash"

export class Client {
    private readonly handlers: {};
    private ws: WebSocket;

    constructor(url :string = "") {
        if (!url) {
            url = ((window.location.protocol === "https:") ? "wss://" : "ws://") + window.location.host + "/ws"
        }
        this.ws = new WebSocket(url, null);
        this.ws.onerror = this.onerror
        this.ws.onmessage = this.onmessage
        this.ws.onclose = this.onclose
        this.ws.onopen = this.onopen
        this.handlers = {}
    }

    public register = (event: wsEvent, handler: any) => {
        if (typeof handler !== 'function') {
            console.log("Invalid handler type: ")
            return
        }
        if (!has(this.handlers, event)) {
            this.handlers[event] = [];
        }
        this.handlers[event].push(handler)
    };

    public emit = (event: wsEvent, payload: object) => {
        if (!has(this.handlers, event)) {
            console.log(`Ignored emit handler: ${wsEvent}`)
            return
        }
        forEach(this.handlers[event], (fn) => {
            fn(payload)
        })
    };

    public onopen = (event: Event) => {
        console.log(`Connected to ${this.ws.url}`)
    };

    public onclose = (event: CloseEvent) => {
        console.log(`Connected to ${this.ws.url}`)
    };

    public  onerror = (event: Event) => {
        console.log(`WS ERR: ${event}`)
    };

    public onmessage = (event: MessageEvent) => {
        const m = <WSEventPayload>JSON.parse(event.data);
        this.emit(m.wsEvent, m.payload)
    };

    public send = (payload: WSEventPayload) => {
        try {
            this.ws.send(payload.encode())
        } catch (e) {
            console.log(`Failed to send payload ${payload.wsEvent}`)
        }
    };

    public open_mission = (mission_id: number) => {
        const req = new WSEventPayload(wsEvent.missionOpen, {"mission_id": mission_id});
        this.send(req)
    };

    public message_send = (msg: string) => {
        this.send(new WSEventPayload(wsEvent.missionSendMessage, {'message': msg}))
    };
}

export enum wsEvent {
    missionOpen = 1,
    missionEvents = 2,
    missionSendMessage = 3,
    missionRecvMessage = 4,
    missionNewFlight = 5,
    missionPosition = 6
}

class WSEventPayload {
    public wsEvent: wsEvent
    public payload: object

    constructor(event: wsEvent, payload: any) {
        this.wsEvent = event
        this.payload = payload
    }

    encode(): string {
        try {
            return JSON.stringify(this);
        } catch (e) {
            console.log(`Failed to encode payload: ${this}`,);
            throw e;
        }
    }
}

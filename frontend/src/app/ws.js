
export const EvtConnect = 1;
export const EvtPing = 2
export const EvtPong = 3;
export const EvtMessage = 10;
export const EvtSetMission = 20
export const EvtError = 10000

function rand_string(length) {
    let result           = '';
    const characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    const charactersLength = characters.length;
    for ( let i = 0; i < length; i++ ) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }
    return result;
}

export class WSClient {
    constructor() {
        this._open = false;
        this.sent_ping = false;
        this._url = `ws://${window.location.host}/ws`
        this._ws = new WebSocket(this._url)
        this._ws.onmessage = this._onmessage.bind(this);
        this._ws.onopen = this._onopen.bind(this);
        this._ws.onclose = this._onclose.bind(this);
        this._ws.onerror = this._onerror.bind(this);
        this._handlers = {};
        this.register(EvtPong, this._pong)
        this.last_ping_data = "";
        this.last_ping_time = null;
    }

    onopen(event) {
        console.log("default onopen")
    };

    _ping() {
        if (this.sent_ping) {
            return
        }
        this.last_ping_data = rand_string(5)
        if (this.send(EvtPing, {'data': this.last_ping_data})) {
            this.last_ping_time = new Date();
            this.sent_ping = true;
            console.log("sent ping")
        }
    }

    _pong(event) {
        if (!this.sent_ping) {
            return
        }
        const {data} = event;
        if (data !== this.last_ping_data) {
            console.log("Got invalid ping data?", data, this.last_ping_data)
        }
        this.last_ping_data = data;
        this.sent_ping = false;
        console.log("Got pong'd");
    }

    register(event, handler) {
        this._handlers[event] = handler;
    }

    _onopen(event) {
        console.log("Connection opened")
        this._open = true
        this.onopen(event)
        this.heartbeat_func = setInterval(this._ping.bind(this), 30000)

    }

    _onclose(event) {
        console.log(`Connection closed: Code=${event.code} Reason=${event.reason} Clean=${event.wasClean}`);
        clearInterval(this.heartbeat_func)
        this.heartbeat_func = null;
    }

    _onerror(event) {
        console.log("Connection error", event.err)
        this._open = false;
    }

    _onmessage(event) {
        const payload = JSON.parse(event.data);
        if (!this._handlers.hasOwnProperty(event)) {
            console.log("Unhandled event payload", event)
            return
        }
        if (!(this._handlers[event] instanceof Function)) {
            console.log("Handler value isn't function", event);
            return
        }
        try {
            this._handlers[event](payload);
        } catch (e) {
            console.log("Error handling ws event", e.message)
        }
    }

    send(event, payload) {
        try {
            this._ws.send(JSON.stringify({
                "event": event,
                "payload": payload,
            }))
            return true
        } catch (e) {
            console.log("Error sending ws event", e.message)
            return false
        }
    }
}
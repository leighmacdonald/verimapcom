class WS {
    private readonly url: string;
    private _ws: WebSocket;

    constructor(url) {
        this.url = url;
        this._ws = new WebSocket(url);
        this._ws.onopen = this._onopen;
        this._ws.onclose = this._onclose;
        this._ws.onerror = this._onerror;
        this._ws.onmessage = this._onmessage;
    }

    private _onopen(event: Event) {
        console.log(`Connected to ${this.url}`)
    }

    private _onclose(event: CloseEvent) {
        console.log(`Connected to ${this.url}`)
    }

    private _onerror(event: Event) {
        console.log(`WS ERR: ${event}`)
    }

    private _onmessage(event: MessageEvent) {
        const m = JSON.parse(event.data);

    }

    public send(payload: WSEventPayload) {
        try {
            this._ws.send(payload.encode())
        } catch (e) {
            console.log(`Failed to send payload ${payload.wsEvent}`)
        }
    }
}

enum wsEvent {
    openMission = 1
}

class WSEventPayload {
    public wsEvent: wsEvent
    public payload: object

    constructor(event: wsEvent, payload: object) {
        this.wsEvent = event
        this.payload = payload
    }

    encode() : string {
        try {
            return JSON.stringify(this);
        } catch (e) {
            console.log(`Failed to encode payload: ${this}`, );
            throw e;
        }
    }
}


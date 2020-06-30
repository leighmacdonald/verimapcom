/**
 *
 * @param client
 * @param mission_id
 */
export function open_mission(client: RPCClient, mission_id: number) {
    const req = new MissionRequest();
    req.setMissionId(mission_id);
    const call = client.syncOpenMission(req, {}, (err: grpcWeb.Error, response: MissionReply) => {
        console.log(response.getMessage());
    })
    call.on('status', (status: grpcWeb.Status) => {
        console.log("got call on status");
    })
}
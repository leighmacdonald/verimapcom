import * as jspb from "google-protobuf"

import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as google_protobuf_any_pb from 'google-protobuf/google/protobuf/any_pb';

export class CreateFlightRequest extends jspb.Message {
  getDescription(): string;
  setDescription(value: string): CreateFlightRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateFlightRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateFlightRequest): CreateFlightRequest.AsObject;
  static serializeBinaryToWriter(message: CreateFlightRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateFlightRequest;
  static deserializeBinaryFromReader(message: CreateFlightRequest, reader: jspb.BinaryReader): CreateFlightRequest;
}

export namespace CreateFlightRequest {
  export type AsObject = {
    description: string,
  }
}

export class CreateFlightResponse extends jspb.Message {
  getFlightId(): number;
  setFlightId(value: number): CreateFlightResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateFlightResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateFlightResponse): CreateFlightResponse.AsObject;
  static serializeBinaryToWriter(message: CreateFlightResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateFlightResponse;
  static deserializeBinaryFromReader(message: CreateFlightResponse, reader: jspb.BinaryReader): CreateFlightResponse;
}

export namespace CreateFlightResponse {
  export type AsObject = {
    flightId: number,
  }
}

export class MissionReply extends jspb.Message {
  getMissionId(): number;
  setMissionId(value: number): MissionReply;

  getMessage(): string;
  setMessage(value: string): MissionReply;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MissionReply.AsObject;
  static toObject(includeInstance: boolean, msg: MissionReply): MissionReply.AsObject;
  static serializeBinaryToWriter(message: MissionReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MissionReply;
  static deserializeBinaryFromReader(message: MissionReply, reader: jspb.BinaryReader): MissionReply;
}

export namespace MissionReply {
  export type AsObject = {
    missionId: number,
    message: string,
  }
}

export class PingRequest extends jspb.Message {
  getAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setAt(value?: google_protobuf_timestamp_pb.Timestamp): PingRequest;
  hasAt(): boolean;
  clearAt(): PingRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PingRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PingRequest): PingRequest.AsObject;
  static serializeBinaryToWriter(message: PingRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PingRequest;
  static deserializeBinaryFromReader(message: PingRequest, reader: jspb.BinaryReader): PingRequest;
}

export namespace PingRequest {
  export type AsObject = {
    at?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class FileUpload extends jspb.Message {
  getFileName(): string;
  setFileName(value: string): FileUpload;

  getFileSize(): number;
  setFileSize(value: number): FileUpload;

  getData(): Uint8Array | string;
  getData_asU8(): Uint8Array;
  getData_asB64(): string;
  setData(value: Uint8Array | string): FileUpload;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FileUpload.AsObject;
  static toObject(includeInstance: boolean, msg: FileUpload): FileUpload.AsObject;
  static serializeBinaryToWriter(message: FileUpload, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FileUpload;
  static deserializeBinaryFromReader(message: FileUpload, reader: jspb.BinaryReader): FileUpload;
}

export namespace FileUpload {
  export type AsObject = {
    fileName: string,
    fileSize: number,
    data: Uint8Array | string,
  }
}

export class FileReply extends jspb.Message {
  getStatus(): number;
  setStatus(value: number): FileReply;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FileReply.AsObject;
  static toObject(includeInstance: boolean, msg: FileReply): FileReply.AsObject;
  static serializeBinaryToWriter(message: FileReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FileReply;
  static deserializeBinaryFromReader(message: FileReply, reader: jspb.BinaryReader): FileReply;
}

export namespace FileReply {
  export type AsObject = {
    status: number,
  }
}

export class Location extends jspb.Message {
  getLat(): number;
  setLat(value: number): Location;

  getLon(): number;
  setLon(value: number): Location;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Location.AsObject;
  static toObject(includeInstance: boolean, msg: Location): Location.AsObject;
  static serializeBinaryToWriter(message: Location, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Location;
  static deserializeBinaryFromReader(message: Location, reader: jspb.BinaryReader): Location;
}

export namespace Location {
  export type AsObject = {
    lat: number,
    lon: number,
  }
}

export class PositionEvent extends jspb.Message {
  getAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setAt(value?: google_protobuf_timestamp_pb.Timestamp): PositionEvent;
  hasAt(): boolean;
  clearAt(): PositionEvent;

  getId(): number;
  setId(value: number): PositionEvent;

  getLocation(): Location | undefined;
  setLocation(value?: Location): PositionEvent;
  hasLocation(): boolean;
  clearLocation(): PositionEvent;

  getElevation(): number;
  setElevation(value: number): PositionEvent;

  getFlightId(): number;
  setFlightId(value: number): PositionEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PositionEvent.AsObject;
  static toObject(includeInstance: boolean, msg: PositionEvent): PositionEvent.AsObject;
  static serializeBinaryToWriter(message: PositionEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PositionEvent;
  static deserializeBinaryFromReader(message: PositionEvent, reader: jspb.BinaryReader): PositionEvent;
}

export namespace PositionEvent {
  export type AsObject = {
    at?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    id: number,
    location?: Location.AsObject,
    elevation: number,
    flightId: number,
  }
}

export class HotSpotEvent extends jspb.Message {
  getAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setAt(value?: google_protobuf_timestamp_pb.Timestamp): HotSpotEvent;
  hasAt(): boolean;
  clearAt(): HotSpotEvent;

  getId(): number;
  setId(value: number): HotSpotEvent;

  getLocation(): Location | undefined;
  setLocation(value?: Location): HotSpotEvent;
  hasLocation(): boolean;
  clearLocation(): HotSpotEvent;

  getDelta(): number;
  setDelta(value: number): HotSpotEvent;

  getFlightId(): number;
  setFlightId(value: number): HotSpotEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HotSpotEvent.AsObject;
  static toObject(includeInstance: boolean, msg: HotSpotEvent): HotSpotEvent.AsObject;
  static serializeBinaryToWriter(message: HotSpotEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HotSpotEvent;
  static deserializeBinaryFromReader(message: HotSpotEvent, reader: jspb.BinaryReader): HotSpotEvent;
}

export namespace HotSpotEvent {
  export type AsObject = {
    at?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    id: number,
    location?: Location.AsObject,
    delta: number,
    flightId: number,
  }
}

export class StatusReply extends jspb.Message {
  getStatus(): StatusReply.Status;
  setStatus(value: StatusReply.Status): StatusReply;

  getMessage(): string;
  setMessage(value: string): StatusReply;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StatusReply.AsObject;
  static toObject(includeInstance: boolean, msg: StatusReply): StatusReply.AsObject;
  static serializeBinaryToWriter(message: StatusReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StatusReply;
  static deserializeBinaryFromReader(message: StatusReply, reader: jspb.BinaryReader): StatusReply;
}

export namespace StatusReply {
  export type AsObject = {
    status: StatusReply.Status,
    message: string,
  }

  export enum Status { 
    OK = 0,
    ERROR = 1,
  }
}

export class ChatMessageRequest extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): ChatMessageRequest;

  getAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setAt(value?: google_protobuf_timestamp_pb.Timestamp): ChatMessageRequest;
  hasAt(): boolean;
  clearAt(): ChatMessageRequest;

  getStartIdx(): number;
  setStartIdx(value: number): ChatMessageRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChatMessageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ChatMessageRequest): ChatMessageRequest.AsObject;
  static serializeBinaryToWriter(message: ChatMessageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChatMessageRequest;
  static deserializeBinaryFromReader(message: ChatMessageRequest, reader: jspb.BinaryReader): ChatMessageRequest;
}

export namespace ChatMessageRequest {
  export type AsObject = {
    message: string,
    at?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    startIdx: number,
  }
}

export class MissionEvent extends jspb.Message {
  getEventType(): EventType;
  setEventType(value: EventType): MissionEvent;

  getPayload(): google_protobuf_any_pb.Any | undefined;
  setPayload(value?: google_protobuf_any_pb.Any): MissionEvent;
  hasPayload(): boolean;
  clearPayload(): MissionEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MissionEvent.AsObject;
  static toObject(includeInstance: boolean, msg: MissionEvent): MissionEvent.AsObject;
  static serializeBinaryToWriter(message: MissionEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MissionEvent;
  static deserializeBinaryFromReader(message: MissionEvent, reader: jspb.BinaryReader): MissionEvent;
}

export namespace MissionEvent {
  export type AsObject = {
    eventType: EventType,
    payload?: google_protobuf_any_pb.Any.AsObject,
  }
}

export class MissionRequest extends jspb.Message {
  getMissionId(): number;
  setMissionId(value: number): MissionRequest;

  getStartIdx(): number;
  setStartIdx(value: number): MissionRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MissionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: MissionRequest): MissionRequest.AsObject;
  static serializeBinaryToWriter(message: MissionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MissionRequest;
  static deserializeBinaryFromReader(message: MissionRequest, reader: jspb.BinaryReader): MissionRequest;
}

export namespace MissionRequest {
  export type AsObject = {
    missionId: number,
    startIdx: number,
  }
}

export class CreateMissionRequest extends jspb.Message {
  getName(): number;
  setName(value: number): CreateMissionRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateMissionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateMissionRequest): CreateMissionRequest.AsObject;
  static serializeBinaryToWriter(message: CreateMissionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateMissionRequest;
  static deserializeBinaryFromReader(message: CreateMissionRequest, reader: jspb.BinaryReader): CreateMissionRequest;
}

export namespace CreateMissionRequest {
  export type AsObject = {
    name: number,
  }
}

export enum EventType { 
  EVT_UNKNOWN = 0,
  EVT_CONNECT = 1,
  EVT_PING = 2,
  EVT_PONG = 3,
  EVT_MESSAGE = 10,
  EVT_SET_MISSION = 20,
  EVT_ERROR = 1000,
}

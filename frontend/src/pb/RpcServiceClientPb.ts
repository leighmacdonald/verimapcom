/**
 * @fileoverview gRPC-Web generated client stub for 
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as google_protobuf_any_pb from 'google-protobuf/google/protobuf/any_pb';

import {
  ChatMessageRequest,
  CreateFlightRequest,
  CreateFlightResponse,
  CreateMissionRequest,
  FileReply,
  FileUpload,
  HotSpotEvent,
  MissionEvent,
  MissionReply,
  MissionRequest,
  PositionEvent,
  StatusReply} from './rpc_pb';

export class RPCClient {
  client_: grpcWeb.AbstractClientBase;
  hostname_: string;
  credentials_: null | { [index: string]: string; };
  options_: null | { [index: string]: string; };

  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; }) {
    if (!options) options = {};
    if (!credentials) credentials = {};
    options['format'] = 'text';

    this.client_ = new grpcWeb.GrpcWebClientBase(options);
    this.hostname_ = hostname;
    this.credentials_ = credentials;
    this.options_ = options;
  }

  methodInfoClientStreamMissionEvents = new grpcWeb.AbstractClientBase.MethodInfo(
    MissionEvent,
    (request: MissionRequest) => {
      return request.serializeBinary();
    },
    MissionEvent.deserializeBinary
  );

  clientStreamMissionEvents(
    request: MissionRequest,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      new URL('/RPC/ClientStreamMissionEvents', this.hostname_).toString(),
      request,
      metadata || {},
      this.methodInfoClientStreamMissionEvents);
  }

  methodInfoClientStreamPositions = new grpcWeb.AbstractClientBase.MethodInfo(
    PositionEvent,
    (request: MissionRequest) => {
      return request.serializeBinary();
    },
    PositionEvent.deserializeBinary
  );

  clientStreamPositions(
    request: MissionRequest,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      new URL('/RPC/ClientStreamPositions', this.hostname_).toString(),
      request,
      metadata || {},
      this.methodInfoClientStreamPositions);
  }

  methodInfoClientStreamHotSpots = new grpcWeb.AbstractClientBase.MethodInfo(
    HotSpotEvent,
    (request: MissionRequest) => {
      return request.serializeBinary();
    },
    HotSpotEvent.deserializeBinary
  );

  clientStreamHotSpots(
    request: MissionRequest,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      new URL('/RPC/ClientStreamHotSpots', this.hostname_).toString(),
      request,
      metadata || {},
      this.methodInfoClientStreamHotSpots);
  }

  methodInfoClientSendMessage = new grpcWeb.AbstractClientBase.MethodInfo(
    StatusReply,
    (request: ChatMessageRequest) => {
      return request.serializeBinary();
    },
    StatusReply.deserializeBinary
  );

  clientSendMessage(
    request: ChatMessageRequest,
    metadata: grpcWeb.Metadata | null): Promise<StatusReply>;

  clientSendMessage(
    request: ChatMessageRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: StatusReply) => void): grpcWeb.ClientReadableStream<StatusReply>;

  clientSendMessage(
    request: ChatMessageRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: StatusReply) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        new URL('/RPC/ClientSendMessage', this.hostname_).toString(),
        request,
        metadata || {},
        this.methodInfoClientSendMessage,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/RPC/ClientSendMessage',
    request,
    metadata || {},
    this.methodInfoClientSendMessage);
  }

  methodInfoSyncSendFile = new grpcWeb.AbstractClientBase.MethodInfo(
    FileReply,
    (request: FileUpload) => {
      return request.serializeBinary();
    },
    FileReply.deserializeBinary
  );

  syncSendFile(
    request: FileUpload,
    metadata: grpcWeb.Metadata | null): Promise<FileReply>;

  syncSendFile(
    request: FileUpload,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: FileReply) => void): grpcWeb.ClientReadableStream<FileReply>;

  syncSendFile(
    request: FileUpload,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: FileReply) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        new URL('/RPC/SyncSendFile', this.hostname_).toString(),
        request,
        metadata || {},
        this.methodInfoSyncSendFile,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/RPC/SyncSendFile',
    request,
    metadata || {},
    this.methodInfoSyncSendFile);
  }

  methodInfoSyncCreateFlight = new grpcWeb.AbstractClientBase.MethodInfo(
    CreateFlightResponse,
    (request: CreateFlightRequest) => {
      return request.serializeBinary();
    },
    CreateFlightResponse.deserializeBinary
  );

  syncCreateFlight(
    request: CreateFlightRequest,
    metadata: grpcWeb.Metadata | null): Promise<CreateFlightResponse>;

  syncCreateFlight(
    request: CreateFlightRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: CreateFlightResponse) => void): grpcWeb.ClientReadableStream<CreateFlightResponse>;

  syncCreateFlight(
    request: CreateFlightRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: CreateFlightResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        new URL('/RPC/SyncCreateFlight', this.hostname_).toString(),
        request,
        metadata || {},
        this.methodInfoSyncCreateFlight,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/RPC/SyncCreateFlight',
    request,
    metadata || {},
    this.methodInfoSyncCreateFlight);
  }

  methodInfoSyncCreateMission = new grpcWeb.AbstractClientBase.MethodInfo(
    MissionReply,
    (request: CreateMissionRequest) => {
      return request.serializeBinary();
    },
    MissionReply.deserializeBinary
  );

  syncCreateMission(
    request: CreateMissionRequest,
    metadata: grpcWeb.Metadata | null): Promise<MissionReply>;

  syncCreateMission(
    request: CreateMissionRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: MissionReply) => void): grpcWeb.ClientReadableStream<MissionReply>;

  syncCreateMission(
    request: CreateMissionRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: MissionReply) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        new URL('/RPC/SyncCreateMission', this.hostname_).toString(),
        request,
        metadata || {},
        this.methodInfoSyncCreateMission,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/RPC/SyncCreateMission',
    request,
    metadata || {},
    this.methodInfoSyncCreateMission);
  }

  methodInfoSyncOpenMission = new grpcWeb.AbstractClientBase.MethodInfo(
    MissionReply,
    (request: MissionRequest) => {
      return request.serializeBinary();
    },
    MissionReply.deserializeBinary
  );

  syncOpenMission(
    request: MissionRequest,
    metadata: grpcWeb.Metadata | null): Promise<MissionReply>;

  syncOpenMission(
    request: MissionRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: MissionReply) => void): grpcWeb.ClientReadableStream<MissionReply>;

  syncOpenMission(
    request: MissionRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: MissionReply) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        new URL('/RPC/SyncOpenMission', this.hostname_).toString(),
        request,
        metadata || {},
        this.methodInfoSyncOpenMission,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/RPC/SyncOpenMission',
    request,
    metadata || {},
    this.methodInfoSyncOpenMission);
  }

}


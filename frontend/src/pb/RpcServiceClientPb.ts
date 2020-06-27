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

  methodInfoSourceSendFile = new grpcWeb.AbstractClientBase.MethodInfo(
    FileReply,
    (request: FileUpload) => {
      return request.serializeBinary();
    },
    FileReply.deserializeBinary
  );

  sourceSendFile(
    request: FileUpload,
    metadata: grpcWeb.Metadata | null): Promise<FileReply>;

  sourceSendFile(
    request: FileUpload,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: FileReply) => void): grpcWeb.ClientReadableStream<FileReply>;

  sourceSendFile(
    request: FileUpload,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: FileReply) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        new URL('/RPC/SourceSendFile', this.hostname_).toString(),
        request,
        metadata || {},
        this.methodInfoSourceSendFile,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/RPC/SourceSendFile',
    request,
    metadata || {},
    this.methodInfoSourceSendFile);
  }

  methodInfoSendMessage = new grpcWeb.AbstractClientBase.MethodInfo(
    StatusReply,
    (request: ChatMessageRequest) => {
      return request.serializeBinary();
    },
    StatusReply.deserializeBinary
  );

  sendMessage(
    request: ChatMessageRequest,
    metadata: grpcWeb.Metadata | null): Promise<StatusReply>;

  sendMessage(
    request: ChatMessageRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: StatusReply) => void): grpcWeb.ClientReadableStream<StatusReply>;

  sendMessage(
    request: ChatMessageRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: StatusReply) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        new URL('/RPC/SendMessage', this.hostname_).toString(),
        request,
        metadata || {},
        this.methodInfoSendMessage,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/RPC/SendMessage',
    request,
    metadata || {},
    this.methodInfoSendMessage);
  }

  methodInfoCreateFlight = new grpcWeb.AbstractClientBase.MethodInfo(
    CreateFlightResponse,
    (request: CreateFlightRequest) => {
      return request.serializeBinary();
    },
    CreateFlightResponse.deserializeBinary
  );

  createFlight(
    request: CreateFlightRequest,
    metadata: grpcWeb.Metadata | null): Promise<CreateFlightResponse>;

  createFlight(
    request: CreateFlightRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: CreateFlightResponse) => void): grpcWeb.ClientReadableStream<CreateFlightResponse>;

  createFlight(
    request: CreateFlightRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: CreateFlightResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        new URL('/RPC/CreateFlight', this.hostname_).toString(),
        request,
        metadata || {},
        this.methodInfoCreateFlight,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/RPC/CreateFlight',
    request,
    metadata || {},
    this.methodInfoCreateFlight);
  }

  methodInfoCreateMission = new grpcWeb.AbstractClientBase.MethodInfo(
    MissionReply,
    (request: CreateMissionRequest) => {
      return request.serializeBinary();
    },
    MissionReply.deserializeBinary
  );

  createMission(
    request: CreateMissionRequest,
    metadata: grpcWeb.Metadata | null): Promise<MissionReply>;

  createMission(
    request: CreateMissionRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: MissionReply) => void): grpcWeb.ClientReadableStream<MissionReply>;

  createMission(
    request: CreateMissionRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: MissionReply) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        new URL('/RPC/CreateMission', this.hostname_).toString(),
        request,
        metadata || {},
        this.methodInfoCreateMission,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/RPC/CreateMission',
    request,
    metadata || {},
    this.methodInfoCreateMission);
  }

  methodInfoOpenMission = new grpcWeb.AbstractClientBase.MethodInfo(
    MissionReply,
    (request: MissionRequest) => {
      return request.serializeBinary();
    },
    MissionReply.deserializeBinary
  );

  openMission(
    request: MissionRequest,
    metadata: grpcWeb.Metadata | null): Promise<MissionReply>;

  openMission(
    request: MissionRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: MissionReply) => void): grpcWeb.ClientReadableStream<MissionReply>;

  openMission(
    request: MissionRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: MissionReply) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        new URL('/RPC/OpenMission', this.hostname_).toString(),
        request,
        metadata || {},
        this.methodInfoOpenMission,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/RPC/OpenMission',
    request,
    metadata || {},
    this.methodInfoOpenMission);
  }

}


/**
 * @fileoverview gRPC-Web generated client stub for 
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')

var google_protobuf_any_pb = require('google-protobuf/google/protobuf/any_pb.js')
const proto = require('./rpc_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.RPCClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.RPCPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.MissionRequest,
 *   !proto.MissionEvent>}
 */
const methodDescriptor_RPC_ClientStreamMissionEvents = new grpc.web.MethodDescriptor(
  '/RPC/ClientStreamMissionEvents',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.MissionRequest,
  proto.MissionEvent,
  /**
   * @param {!proto.MissionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.MissionEvent.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.MissionRequest,
 *   !proto.MissionEvent>}
 */
const methodInfo_RPC_ClientStreamMissionEvents = new grpc.web.AbstractClientBase.MethodInfo(
  proto.MissionEvent,
  /**
   * @param {!proto.MissionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.MissionEvent.deserializeBinary
);


/**
 * @param {!proto.MissionRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.MissionEvent>}
 *     The XHR Node Readable Stream
 */
proto.RPCClient.prototype.clientStreamMissionEvents =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/RPC/ClientStreamMissionEvents',
      request,
      metadata || {},
      methodDescriptor_RPC_ClientStreamMissionEvents);
};


/**
 * @param {!proto.MissionRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.MissionEvent>}
 *     The XHR Node Readable Stream
 */
proto.RPCPromiseClient.prototype.clientStreamMissionEvents =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/RPC/ClientStreamMissionEvents',
      request,
      metadata || {},
      methodDescriptor_RPC_ClientStreamMissionEvents);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.MissionRequest,
 *   !proto.PositionEvent>}
 */
const methodDescriptor_RPC_ClientStreamPositions = new grpc.web.MethodDescriptor(
  '/RPC/ClientStreamPositions',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.MissionRequest,
  proto.PositionEvent,
  /**
   * @param {!proto.MissionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.PositionEvent.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.MissionRequest,
 *   !proto.PositionEvent>}
 */
const methodInfo_RPC_ClientStreamPositions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.PositionEvent,
  /**
   * @param {!proto.MissionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.PositionEvent.deserializeBinary
);


/**
 * @param {!proto.MissionRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.PositionEvent>}
 *     The XHR Node Readable Stream
 */
proto.RPCClient.prototype.clientStreamPositions =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/RPC/ClientStreamPositions',
      request,
      metadata || {},
      methodDescriptor_RPC_ClientStreamPositions);
};


/**
 * @param {!proto.MissionRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.PositionEvent>}
 *     The XHR Node Readable Stream
 */
proto.RPCPromiseClient.prototype.clientStreamPositions =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/RPC/ClientStreamPositions',
      request,
      metadata || {},
      methodDescriptor_RPC_ClientStreamPositions);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.MissionRequest,
 *   !proto.HotSpotEvent>}
 */
const methodDescriptor_RPC_ClientStreamHotSpots = new grpc.web.MethodDescriptor(
  '/RPC/ClientStreamHotSpots',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.MissionRequest,
  proto.HotSpotEvent,
  /**
   * @param {!proto.MissionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.HotSpotEvent.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.MissionRequest,
 *   !proto.HotSpotEvent>}
 */
const methodInfo_RPC_ClientStreamHotSpots = new grpc.web.AbstractClientBase.MethodInfo(
  proto.HotSpotEvent,
  /**
   * @param {!proto.MissionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.HotSpotEvent.deserializeBinary
);


/**
 * @param {!proto.MissionRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.HotSpotEvent>}
 *     The XHR Node Readable Stream
 */
proto.RPCClient.prototype.clientStreamHotSpots =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/RPC/ClientStreamHotSpots',
      request,
      metadata || {},
      methodDescriptor_RPC_ClientStreamHotSpots);
};


/**
 * @param {!proto.MissionRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.HotSpotEvent>}
 *     The XHR Node Readable Stream
 */
proto.RPCPromiseClient.prototype.clientStreamHotSpots =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/RPC/ClientStreamHotSpots',
      request,
      metadata || {},
      methodDescriptor_RPC_ClientStreamHotSpots);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ChatMessageRequest,
 *   !proto.StatusReply>}
 */
const methodDescriptor_RPC_ClientSendMessage = new grpc.web.MethodDescriptor(
  '/RPC/ClientSendMessage',
  grpc.web.MethodType.UNARY,
  proto.ChatMessageRequest,
  proto.StatusReply,
  /**
   * @param {!proto.ChatMessageRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.StatusReply.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ChatMessageRequest,
 *   !proto.StatusReply>}
 */
const methodInfo_RPC_ClientSendMessage = new grpc.web.AbstractClientBase.MethodInfo(
  proto.StatusReply,
  /**
   * @param {!proto.ChatMessageRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.StatusReply.deserializeBinary
);


/**
 * @param {!proto.ChatMessageRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.StatusReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.StatusReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.RPCClient.prototype.clientSendMessage =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/RPC/ClientSendMessage',
      request,
      metadata || {},
      methodDescriptor_RPC_ClientSendMessage,
      callback);
};


/**
 * @param {!proto.ChatMessageRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.StatusReply>}
 *     A native promise that resolves to the response
 */
proto.RPCPromiseClient.prototype.clientSendMessage =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/RPC/ClientSendMessage',
      request,
      metadata || {},
      methodDescriptor_RPC_ClientSendMessage);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.FileUpload,
 *   !proto.FileReply>}
 */
const methodDescriptor_RPC_SyncSendFile = new grpc.web.MethodDescriptor(
  '/RPC/SyncSendFile',
  grpc.web.MethodType.UNARY,
  proto.FileUpload,
  proto.FileReply,
  /**
   * @param {!proto.FileUpload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.FileReply.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.FileUpload,
 *   !proto.FileReply>}
 */
const methodInfo_RPC_SyncSendFile = new grpc.web.AbstractClientBase.MethodInfo(
  proto.FileReply,
  /**
   * @param {!proto.FileUpload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.FileReply.deserializeBinary
);


/**
 * @param {!proto.FileUpload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.FileReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.FileReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.RPCClient.prototype.syncSendFile =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/RPC/SyncSendFile',
      request,
      metadata || {},
      methodDescriptor_RPC_SyncSendFile,
      callback);
};


/**
 * @param {!proto.FileUpload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.FileReply>}
 *     A native promise that resolves to the response
 */
proto.RPCPromiseClient.prototype.syncSendFile =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/RPC/SyncSendFile',
      request,
      metadata || {},
      methodDescriptor_RPC_SyncSendFile);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.CreateFlightRequest,
 *   !proto.CreateFlightResponse>}
 */
const methodDescriptor_RPC_SyncCreateFlight = new grpc.web.MethodDescriptor(
  '/RPC/SyncCreateFlight',
  grpc.web.MethodType.UNARY,
  proto.CreateFlightRequest,
  proto.CreateFlightResponse,
  /**
   * @param {!proto.CreateFlightRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.CreateFlightResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.CreateFlightRequest,
 *   !proto.CreateFlightResponse>}
 */
const methodInfo_RPC_SyncCreateFlight = new grpc.web.AbstractClientBase.MethodInfo(
  proto.CreateFlightResponse,
  /**
   * @param {!proto.CreateFlightRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.CreateFlightResponse.deserializeBinary
);


/**
 * @param {!proto.CreateFlightRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.CreateFlightResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.CreateFlightResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.RPCClient.prototype.syncCreateFlight =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/RPC/SyncCreateFlight',
      request,
      metadata || {},
      methodDescriptor_RPC_SyncCreateFlight,
      callback);
};


/**
 * @param {!proto.CreateFlightRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.CreateFlightResponse>}
 *     A native promise that resolves to the response
 */
proto.RPCPromiseClient.prototype.syncCreateFlight =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/RPC/SyncCreateFlight',
      request,
      metadata || {},
      methodDescriptor_RPC_SyncCreateFlight);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.CreateMissionRequest,
 *   !proto.MissionReply>}
 */
const methodDescriptor_RPC_SyncCreateMission = new grpc.web.MethodDescriptor(
  '/RPC/SyncCreateMission',
  grpc.web.MethodType.UNARY,
  proto.CreateMissionRequest,
  proto.MissionReply,
  /**
   * @param {!proto.CreateMissionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.MissionReply.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.CreateMissionRequest,
 *   !proto.MissionReply>}
 */
const methodInfo_RPC_SyncCreateMission = new grpc.web.AbstractClientBase.MethodInfo(
  proto.MissionReply,
  /**
   * @param {!proto.CreateMissionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.MissionReply.deserializeBinary
);


/**
 * @param {!proto.CreateMissionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.MissionReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.MissionReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.RPCClient.prototype.syncCreateMission =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/RPC/SyncCreateMission',
      request,
      metadata || {},
      methodDescriptor_RPC_SyncCreateMission,
      callback);
};


/**
 * @param {!proto.CreateMissionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.MissionReply>}
 *     A native promise that resolves to the response
 */
proto.RPCPromiseClient.prototype.syncCreateMission =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/RPC/SyncCreateMission',
      request,
      metadata || {},
      methodDescriptor_RPC_SyncCreateMission);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.MissionRequest,
 *   !proto.MissionReply>}
 */
const methodDescriptor_RPC_SyncOpenMission = new grpc.web.MethodDescriptor(
  '/RPC/SyncOpenMission',
  grpc.web.MethodType.UNARY,
  proto.MissionRequest,
  proto.MissionReply,
  /**
   * @param {!proto.MissionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.MissionReply.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.MissionRequest,
 *   !proto.MissionReply>}
 */
const methodInfo_RPC_SyncOpenMission = new grpc.web.AbstractClientBase.MethodInfo(
  proto.MissionReply,
  /**
   * @param {!proto.MissionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.MissionReply.deserializeBinary
);


/**
 * @param {!proto.MissionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.MissionReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.MissionReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.RPCClient.prototype.syncOpenMission =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/RPC/SyncOpenMission',
      request,
      metadata || {},
      methodDescriptor_RPC_SyncOpenMission,
      callback);
};


/**
 * @param {!proto.MissionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.MissionReply>}
 *     A native promise that resolves to the response
 */
proto.RPCPromiseClient.prototype.syncOpenMission =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/RPC/SyncOpenMission',
      request,
      metadata || {},
      methodDescriptor_RPC_SyncOpenMission);
};


module.exports = proto;


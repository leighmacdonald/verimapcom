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
const methodDescriptor_RPC_StreamMissionEvents = new grpc.web.MethodDescriptor(
  '/RPC/StreamMissionEvents',
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
const methodInfo_RPC_StreamMissionEvents = new grpc.web.AbstractClientBase.MethodInfo(
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
proto.RPCClient.prototype.streamMissionEvents =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/RPC/StreamMissionEvents',
      request,
      metadata || {},
      methodDescriptor_RPC_StreamMissionEvents);
};


/**
 * @param {!proto.MissionRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.MissionEvent>}
 *     The XHR Node Readable Stream
 */
proto.RPCPromiseClient.prototype.streamMissionEvents =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/RPC/StreamMissionEvents',
      request,
      metadata || {},
      methodDescriptor_RPC_StreamMissionEvents);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.MissionRequest,
 *   !proto.PositionEvent>}
 */
const methodDescriptor_RPC_StreamPositions = new grpc.web.MethodDescriptor(
  '/RPC/StreamPositions',
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
const methodInfo_RPC_StreamPositions = new grpc.web.AbstractClientBase.MethodInfo(
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
proto.RPCClient.prototype.streamPositions =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/RPC/StreamPositions',
      request,
      metadata || {},
      methodDescriptor_RPC_StreamPositions);
};


/**
 * @param {!proto.MissionRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.PositionEvent>}
 *     The XHR Node Readable Stream
 */
proto.RPCPromiseClient.prototype.streamPositions =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/RPC/StreamPositions',
      request,
      metadata || {},
      methodDescriptor_RPC_StreamPositions);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.MissionRequest,
 *   !proto.HotSpotEvent>}
 */
const methodDescriptor_RPC_StreamHotSpots = new grpc.web.MethodDescriptor(
  '/RPC/StreamHotSpots',
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
const methodInfo_RPC_StreamHotSpots = new grpc.web.AbstractClientBase.MethodInfo(
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
proto.RPCClient.prototype.streamHotSpots =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/RPC/StreamHotSpots',
      request,
      metadata || {},
      methodDescriptor_RPC_StreamHotSpots);
};


/**
 * @param {!proto.MissionRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.HotSpotEvent>}
 *     The XHR Node Readable Stream
 */
proto.RPCPromiseClient.prototype.streamHotSpots =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/RPC/StreamHotSpots',
      request,
      metadata || {},
      methodDescriptor_RPC_StreamHotSpots);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ChatMessageRequest,
 *   !proto.StatusReply>}
 */
const methodDescriptor_RPC_SendMessage = new grpc.web.MethodDescriptor(
  '/RPC/SendMessage',
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
const methodInfo_RPC_SendMessage = new grpc.web.AbstractClientBase.MethodInfo(
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
proto.RPCClient.prototype.sendMessage =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/RPC/SendMessage',
      request,
      metadata || {},
      methodDescriptor_RPC_SendMessage,
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
proto.RPCPromiseClient.prototype.sendMessage =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/RPC/SendMessage',
      request,
      metadata || {},
      methodDescriptor_RPC_SendMessage);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.FileUpload,
 *   !proto.FileReply>}
 */
const methodDescriptor_RPC_SendFile = new grpc.web.MethodDescriptor(
  '/RPC/SendFile',
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
const methodInfo_RPC_SendFile = new grpc.web.AbstractClientBase.MethodInfo(
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
proto.RPCClient.prototype.sendFile =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/RPC/SendFile',
      request,
      metadata || {},
      methodDescriptor_RPC_SendFile,
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
proto.RPCPromiseClient.prototype.sendFile =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/RPC/SendFile',
      request,
      metadata || {},
      methodDescriptor_RPC_SendFile);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.PingRequest,
 *   !proto.PingReply>}
 */
const methodDescriptor_RPC_Ping = new grpc.web.MethodDescriptor(
  '/RPC/Ping',
  grpc.web.MethodType.UNARY,
  proto.PingRequest,
  proto.PingReply,
  /**
   * @param {!proto.PingRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.PingReply.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.PingRequest,
 *   !proto.PingReply>}
 */
const methodInfo_RPC_Ping = new grpc.web.AbstractClientBase.MethodInfo(
  proto.PingReply,
  /**
   * @param {!proto.PingRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.PingReply.deserializeBinary
);


/**
 * @param {!proto.PingRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.PingReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.PingReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.RPCClient.prototype.ping =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/RPC/Ping',
      request,
      metadata || {},
      methodDescriptor_RPC_Ping,
      callback);
};


/**
 * @param {!proto.PingRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.PingReply>}
 *     A native promise that resolves to the response
 */
proto.RPCPromiseClient.prototype.ping =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/RPC/Ping',
      request,
      metadata || {},
      methodDescriptor_RPC_Ping);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ProjectRequest,
 *   !proto.ProjectReply>}
 */
const methodDescriptor_RPC_OpenProject = new grpc.web.MethodDescriptor(
  '/RPC/OpenProject',
  grpc.web.MethodType.UNARY,
  proto.ProjectRequest,
  proto.ProjectReply,
  /**
   * @param {!proto.ProjectRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ProjectReply.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ProjectRequest,
 *   !proto.ProjectReply>}
 */
const methodInfo_RPC_OpenProject = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ProjectReply,
  /**
   * @param {!proto.ProjectRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ProjectReply.deserializeBinary
);


/**
 * @param {!proto.ProjectRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ProjectReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ProjectReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.RPCClient.prototype.openProject =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/RPC/OpenProject',
      request,
      metadata || {},
      methodDescriptor_RPC_OpenProject,
      callback);
};


/**
 * @param {!proto.ProjectRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ProjectReply>}
 *     A native promise that resolves to the response
 */
proto.RPCPromiseClient.prototype.openProject =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/RPC/OpenProject',
      request,
      metadata || {},
      methodDescriptor_RPC_OpenProject);
};


module.exports = proto;


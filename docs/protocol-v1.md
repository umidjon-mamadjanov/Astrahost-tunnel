Astra Tunnel Protocol v1

Status: Draft for implementation
Version: 1.0
Transport: WebSocket over TLS (WSS)
Internal server listener: "127.0.0.1:7000"

1. Purpose

Astra Tunnel Protocol v1 defines communication between an Astra Tunnel Client and Astra Tunnel Server.

The protocol is responsible for:

- establishing a tunnel session;
- authenticating a client;
- transporting HTTP requests and responses;
- maintaining connection liveness;
- reporting protocol errors;
- closing tunnel sessions cleanly.

The protocol is transport-independent at the logical level, but v1 uses WebSocket as its transport.

2. Transport

Production:

Client
  │
  │ WSS
  ▼
NGINX :443
  │
  │ WebSocket proxy
  ▼
Astra Server 127.0.0.1:7000

Development:

ws://127.0.0.1:7000/connect

Production:

wss://tunnel.example.com/connect

The Astra Server does not expose port "7000" publicly.

3. Protocol Version

Current protocol version:

1

The version is carried in every packet header.

Unsupported versions MUST result in:

ERROR(PROTOCOL_VERSION_UNSUPPORTED)

4. Packet Model

Every protocol message is represented as a packet.

Logical structure:

+--------------------+
| Header             |
+--------------------+
| Payload             |
+--------------------+

The header contains:

Version
Type
Flags
Request ID
Payload Length

5. Header

v1 header:

+----------------+--------+
| Field          | Size   |
+----------------+--------+
| Version        | 1 byte |
| Type           | 1 byte |
| Flags          | 2 byte |
| Request ID     | 16 byte|
| Payload Length | 4 byte |
+----------------+--------+

Total:

24 bytes

All multi-byte integer fields use:

Big Endian

Version

Unsigned 8-bit integer.

1 = Protocol v1

Type

Unsigned 8-bit packet type.

Flags

Unsigned 16-bit field.

v1 initially defines no mandatory flags.

Unknown flags MUST be ignored unless marked as required by a future protocol revision.

Request ID

16-byte identifier.

Recommended representation:

UUID

The Request ID correlates a request with its corresponding response.

For packets that do not require correlation, the Request ID MUST be all zeroes.

Payload Length

Unsigned 32-bit integer.

Maximum payload size for v1:

16 MiB

A packet exceeding this limit MUST be rejected.

6. Packet Types

v1 packet types:

0x01 CONNECT
0x02 CONNECT_OK

0x10 HTTP_REQUEST
0x11 HTTP_RESPONSE

0x20 PING
0x21 PONG

0x30 ERROR
0x31 CLOSE

Reserved ranges:

0x40-0x7F

for future Astra protocol features.

7. CONNECT

"CONNECT" starts a tunnel session.

Direction:

Client → Server

Payload:

{
  "protocol_version": 1,
  "client_version": "1.0.0",
  "platform": "linux",
  "architecture": "amd64",
  "tunnel_name": "demo",
  "auth": "..."
}

Required fields:

protocol_version
client_version
platform
architecture

"tunnel_name" is optional in the initial implementation.

Authentication data is required when server authentication is enabled.

The server MUST validate:

- protocol version;
- client version format;
- supported platform;
- supported architecture;
- authentication;
- requested tunnel name, if supplied.

8. CONNECT_OK

"CONNECT_OK" confirms successful session establishment.

Direction:

Server → Client

Payload:

{
  "protocol_version": 1,
  "session_id": "...",
  "tunnel_id": "...",
  "heartbeat_interval": 20
}

Required fields:

protocol_version
session_id
tunnel_id
heartbeat_interval

After receiving "CONNECT_OK", the client enters:

READY

state.

9. HTTP_REQUEST

Carries an HTTP request from the Astra Server to the Astra Client.

Direction:

Server → Client

Logical payload:

{
  "method": "GET",
  "path": "/",
  "query": "",
  "headers": {},
  "body": "..."
}

The Request ID in the packet header identifies the HTTP request.

The client forwards the request to the configured local application.

Example:

Astra Client
     │
     ▼
http://127.0.0.1:5000/

10. HTTP_RESPONSE

Carries the local application's response back to the Astra Server.

Direction:

Client → Server

Logical payload:

{
  "status": 200,
  "headers": {},
  "body": "..."
}

The Request ID MUST match the corresponding "HTTP_REQUEST".

Example:

HTTP_REQUEST
Request ID: A

      ↓

HTTP_RESPONSE
Request ID: A

11. PING

Used to verify tunnel liveness.

Direction:

Server → Client

Payload may be empty.

The client MUST respond with "PONG".

12. PONG

Used as the response to "PING".

Direction:

Client → Server

The Request ID SHOULD match the corresponding "PING".

13. ERROR

Reports a protocol or session error.

Payload:

{
  "code": "AUTH_FAILED",
  "message": "authentication failed"
}

Initial error codes:

PROTOCOL_VERSION_UNSUPPORTED
INVALID_PACKET
INVALID_PAYLOAD
PAYLOAD_TOO_LARGE
AUTH_REQUIRED
AUTH_FAILED
TUNNEL_NOT_FOUND
SESSION_NOT_FOUND
REQUEST_TIMEOUT
RATE_LIMITED
INTERNAL_ERROR

The "message" field is informational and MUST NOT be relied upon by clients for program logic.

Clients MUST use the error "code".

14. CLOSE

Requests or indicates graceful tunnel closure.

Payload:

{
  "reason": "client_shutdown"
}

Possible reasons:

client_shutdown
server_shutdown
authentication_failed
protocol_error
timeout
administrative

After CLOSE, the session enters:

CLOSED

state.

15. Session State Machine

Client and server sessions use:

CONNECTED
    │
    ▼
HANDSHAKE
    │
    ├── error ──> CLOSED
    │
    ▼
READY
    │
    ├── timeout ──> CLOSED
    ├── error ────> CLOSED
    └── close ────> CLOSED

The server MUST NOT accept "HTTP_REQUEST" traffic from a session that has not reached "READY".

16. Connection Lifecycle

Normal lifecycle:

Client
  │
  │ WebSocket connect
  ▼
Server
  │
  │ CONNECT
  ▼
Server validates
  │
  │ CONNECT_OK
  ▼
Client READY
  │
  │ HTTP traffic
  ▼
Tunnel active
  │
  │ PING/PONG
  ▼
Tunnel remains active
  │
  │ CLOSE
  ▼
Session closed

17. Request Correlation

Every HTTP request receives a unique Request ID.

Example:

Request ID = UUID

Flow:

Browser
  │
  ▼
Server
  │
  │ HTTP_REQUEST
  │ Request ID: A
  ▼
Client
  │
  ▼
localhost:5000
  │
  ▼
Client
  │
  │ HTTP_RESPONSE
  │ Request ID: A
  ▼
Server
  │
  ▼
Browser

The server MUST NOT deliver a response to a different request.

18. Payload Limits

Default limits:

Maximum packet payload: 16 MiB
Maximum HTTP body:      16 MiB

The implementation MUST reject oversized payloads before allocating excessive memory.

Future versions may introduce streaming/chunked payloads.

19. Heartbeat

Default:

Heartbeat interval: 20 seconds
Connection timeout: 60 seconds

Server sends:

PING

Client responds:

PONG

If the server does not receive a valid response within the timeout period, the session is considered dead.

20. Authentication

v1 uses API-key authentication.

The key is supplied during CONNECT.

The server MUST NOT store plaintext API keys.

The server should store a cryptographic hash of the key.

Authentication failure:

CONNECT
   │
   ▼
AUTH_FAILED
   │
   ▼
CLOSE

21. Reconnection

The client MUST support automatic reconnection.

Recommended backoff:

1s
2s
4s
8s
16s
30s
30s
...

Maximum retry delay:

30 seconds

After reconnecting, the client MUST perform a new CONNECT handshake.

A previous session MUST NOT automatically be considered valid after reconnect.

22. Security Requirements

The implementation MUST:

- validate packet lengths;
- validate protocol versions;
- enforce payload limits;
- authenticate clients;
- avoid plaintext API-key storage;
- use WSS in production;
- apply connection timeouts;
- apply request timeouts;
- limit excessive connections;
- recover from unexpected handler errors;
- avoid leaking sensitive authentication data into logs.

23. Compatibility

Protocol v1 implementations MUST reject unsupported protocol versions explicitly.

Future protocol versions may introduce:

v1 → v2

without changing the meaning of existing v1 packet types.

Unknown packet types MUST result in:

ERROR(INVALID_PACKET)

unless the packet type belongs to a reserved extension mechanism defined by a future protocol version.

24. v1 Implementation Boundary

Included:

CONNECT
CONNECT_OK

HTTP_REQUEST
HTTP_RESPONSE

PING
PONG

ERROR
CLOSE

API Key authentication
Heartbeat
Reconnect
WebSocket transport
HTTP tunneling
WSS deployment

Not included:

TCP tunnel
UDP tunnel
QUIC
HTTP/2 transport
Streaming protocol
Dashboard
Billing
Multi-region
Kubernetes
Mobile GUI
Desktop GUI

25. v1 Success Criteria

Astra Tunnel v1 is considered functional when:

[ ] Client connects through WSS
[ ] CONNECT handshake succeeds
[ ] Session receives READY state
[ ] Tunnel is registered
[ ] Public HTTP request reaches client
[ ] Client forwards request to localhost application
[ ] Local response reaches public client
[ ] Request IDs correlate correctly
[ ] Heartbeat works
[ ] Dead connections are detected
[ ] Client reconnects
[ ] Authentication works
[ ] Invalid authentication is rejected
[ ] Oversized packets are rejected
[ ] Graceful shutdown works
[ ] Dom Cloud deployment works
[ ] go test ./... passes
[ ] go vet ./... passes

26. Protocol Principle

Astra Tunnel v1 follows one primary principle:

Transport carries packets.
Protocol defines meaning.
Engine processes packets.
Session represents a connection.
Registry identifies tunnels.
Proxy handles HTTP traffic.

This separation MUST be preserved throughout v1.0 development.

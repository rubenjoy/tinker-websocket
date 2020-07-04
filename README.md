# Tinkering with Websocket

I'm in the middle of working with websocket and still kinda new.  
I understand what Websocket is in the surface but blindly know nothing 
under the hood. This repo should stash away my tinkering-websocket
journey so far. 

## Objectives

 - faketime chat, that delaying message delivery. The delay time may adjusted
   to zero means it realtime.
 - understand how to use [autobahn](https://github.com/crossbario/autobahn-testsuite)
   test suite, create a report and interpret it.

## Tools and Library

 - [websocat](https://github.com/vi/websocat): more advanced netcat-like
   utility written in rust.
 - [wscat](https://github.com/websockets/wscat): netcat-like utility to listen
   incoming websocket request or connect as client.
 - [websockets](https://pypi.org/project/websockets/): a python library that
   support client/server websocket, readme: [here](https://websockets.readthedocs.io/en/stable/).
 - [gorilla websocket](https://github.com/gorilla/websocket): a golang library that
   upgrade incoming HTTP request into websocket. It supports dialling a websocket too.
 - [websocket/ws](https://github.com/websockets/wscat): a NodeJS websocket
   library.
 - [golang net/websocket](https://godoc.org/golang.org/x/net/websocket)

## References

 - [MDN Web API Websocket](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)


## Notes

The transmitted data message can be text, binary or control message.
There three control messages:
 - /close 1006
 - /ping
 - /pong
 
Whenever a ping request accepted, the websocket entity should response
with pong control message. The ping/pong control message acts as heartbeat 
and probes the healthiness of a connection.





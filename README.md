# Bambus

<img src="https://www.svgrepo.com/show/148390/bamboo.svg" width="100">

----
**UNDER DEVELOPMENT**

Bambus is a real-time scalable messaging server, enables web frontend (typically browsers) to get real time updates from web backends. 

# Introduction

To provide real time user experience, backends need to send updates to UIs in a pubsub fashion. 
So we need kind of pubsub broker that sits between backends and frontends.
Then backends can publish events to channels that frontends are subscribed to.

Fortunately we have a bunch of powerful pubsub brokers.
They are all amazing, but frontends cannot connect directly to them,
we still need another piece that enables connections to web frontends and applies authentication.
This is what bambus tries to do.

It sits in front of a pubsub broker (redis currently supported) and accepts socket connections (Websockets, SSEs) from web frontends.

For each incoming socket connection, bambus performs the authentication process to figure out which channels of events the socket should receive,
then it assigns it to a one of the running threads (go routines).
Each thread handles one open broker connection and multiple frontend sockets.


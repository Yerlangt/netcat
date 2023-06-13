# net-cat
## Authors: Yerlangt
### Objectives

This project consists on recreating the NetCat in a Server-Client Architecture that can run in a server mode on a specified port listening for incoming connections, and it can be used in client mode, trying to connect to a specified port and transmitting information to the server.

### Features of Net-cat

    TCP connection between server and multiple clients (relation of 1 to many).
    A name requirement to the client.
    Control connections quantity.
    Clients must be able to send messages to the chat.
    Do not broadcast EMPTY messages from a client.
    Messages sent, must be identified by the time that was sent and the user name of who sent the message, example : [2020-01-20 15:48:41][client.name]:[client.message]
    If a Client joins the chat, all the previous messages sent to the chat must be uploaded to the new Client.
    If a Client connects to the server, the rest of the Clients must be informed by the server that the Client joined the group.
    If a Client exits the chat, the rest of the Clients must be informed by the server that the Client left.
    All Clients must receive the messages sent by other Clients.
    If a Client leaves the chat, the rest of the Clients must not disconnect.


### Usage 
```
go run cmd/main.go [port]
```
By default port 8989 will be used

It also possible to run the net-cat with the use of the makefile (but only with default port): 
```
make run
```

For global connection (other computers in the same network) please follow next steps:
```
1. In terminal:
   - ifconfig
2. Take the information regarding the host of the eno1 for the inet (host for global)
3. Run the program on localhost with port 27960 (will allow global communication):
   - go run cmd/main.go 27960
4. To connect from the same computer to the chat:
   - nc localhost 27960
5. To connect from another computer to the chat:
   - nc [host fot global] 27960
```


### Additional properties

```
There is possibility to change user's name with command (should be asked as message by user). Name will be changed if it is unique, and fits properties for nickname. Should be entered without spaces:
    --CN:[new nickname]

There is possibility to check the status of the chat - total number of users in the chat at the moment (should be asked as message by user):
    !StatusCheck

To create a binary file use makefile, in terminal box:
    make build
```

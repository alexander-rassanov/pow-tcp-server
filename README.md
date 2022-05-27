# pow-tcp-server

## Challenge

"Word of Wisdom" tcp server. • TCP server should be protected from DDOS attacks with the Prof of Work
(https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
* The choice of the POW algorithm should be explained.
* After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
* Docker file should be provided both for the server and for the client that solves the POW challenge

## Applications

The repository contains two applications:
- [client](cmd/client/) - application interacts with a server to receive quotes of wisdom. It also supports hash cash algorithm in order to get access to the server.
- [server](cmd/server/) - application handles incoming requests and response quotes of wisdom. It requires a computation of hash cash from clients in order to serve requests.

## How it works

Client and server interacts through TCP with a Challenge–response protocol. Client and server sends data using specific packages with
specific headers with payload.
Example of the interaction:

1. Client - Header(Request Challenge) Payload(Empty) -> Server
2. Server - Header(Response Challenge) Payload(Hash cash data) -> Client
3. Client - Header(Request Service) Payload(resolved Hash cash data with correct fields) -> Server
4. Server - Header(Response Service) Payload(A string of quote) -> Client

## Quick start

Makefile provides useful targets for users. It may help to run and test the application quickly.

| Target             | Description                                                                 |
|--------------------|-----------------------------------------------------------------------------|
| run-client         | build and run client using golang on the host where the target is executed. |
| run-server         | build and run server using golang on the host where the target is executed. |
| run-server         | build and run server using golang on the host where the target is executed. |
| lint               | executes lint to make sure code is consistent and clean                     |
| build-client-image | build a Docker image of the client                                          |
| build-server-image | build a Docker image of the server                                          |
| test               | runs tests against all packages                                             |






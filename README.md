# dinosaur

A game about writing code like a dinosaur

## Concept

- User picks a language (presently only Go is supported)
- User is presented w/ 2 panes; editor and live feed to their code being executed
- User makes changes to get code doing what it should do

## Ideas

Not actually sure what the challenges will be yet; I'm thinking multiplayer (e.g sockets) type things; e.g.:

- Two players need to exchange some messages with each other
- Two players need to complete a set of calculations (over a series of messages) and each player can only see the odd or even half of the
  steps
- Two or more players need to develop a decentralised chat system (discovery etc) without knowing who's out there
- Same as above, but now there are adversaries (also players) actively trying to hinder the other players

I might try and include some sort of in-game chat and hope people don't use it to cheat; but maybe there's a hard mode for extra points
where the in-game chat is disabled (so the implication is that step one is discovery of your team mates e.g. using multicast, validation
that they're not adversaries in some way and all team chat is done with code written in the session).

## Architecture

- Backend
    - [Go](https://go.dev/) w/ built-in HTTP server and subprocess orchestration libs
    - [entr](https://github.com/eradman/entr)
    - [sorenisanerd's](https://github.com/sorenisanerd) fork of [gotty](https://github.com/sorenisanerd/gotty)
- Frontend
    - [Create React App](https://create-react-app.dev/)
    - [Microsoft Monaco Editor](https://github.com/microsoft/monaco-editor)

## Flows

- `Frontend` makes a GET request to the backend to be allocated a session
- `Frontend` mounts `Shell` component that gives the live feed of the `gotty` session (just as an iframe)
- `User` makes edits in `Editor` component
- `Editor` component text POST'd to `Backend`
- `entr` in `Backend`  re-runs code process
- `Shell` component in frontend continues to display live feed

## TODO

- Containerise everything (right now it runs natively)
    - Probably Docker Compose w/ private networks (no internet access) to limit exploitability
- Put an Nginx reverse proxy in front of it all
- Add support for other languages (should easy if Monaco supports them and I can work out the MVP env + boilerplate)
    - Python
    - TypeScript
    - C
    - Rust
    - Java
- Ability to relate sessions together in a game
- Recording for all interactions for a game
- All the actual gamification stuff
    - Identity / single sign-on
    - Lobby
    - Scoring

## How to run it

### Prerequisites

- [Go](https://go.dev/)
- [entr](https://github.com/eradman/entr)
- [sorenisanerd's](https://github.com/sorenisanerd) fork of [gotty](https://github.com/sorenisanerd/gotty)
- [npm](https://docs.npmjs.com/)

### Steps

**Backend**

```shell
go run cmd/main.go
```

**Frontend**

```shell
npm ci
npm run start
```

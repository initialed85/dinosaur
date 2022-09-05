# dinosaur

A training tool that gets you to write code like a dinosaur.

## Concept

- User picks a language
- User is presented w/ 2 panes; editor and live feed to their code being executed
- User makes changes to code to achieve some goal

## Ideas

I don't really know what the goals will be yet, but I'm thinking about about multiplayer challenges- something like:

- Two players need to exchange some messages with each other
- Two players need to complete a set of calculations (over a series of messages) and each player can only see the odd or even half of the
  steps
- Two or more players need to develop a decentralised chat system (discovery etc) without knowing who's out there
- Same as above, but now there are adversaries (also players) actively trying to hinder the other players

Ultra hard mode would be players on the same team need to discover and communicate with each other using only the code in front of them; in
the case where there are adversaries there'd be a whole element of validating that the person you're talking to is a teammate, not an
adversary.

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

- Ability to relate sessions together in a game
- Recording for all interactions for a game
- All the actual gamification stuff
    - Identity / single sign-on
    - Lobby
    - Scoring

## How to run it

### Prerequisites

- [Go](https://go.dev/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### Steps

```shell
./run.sh
```

## How to use it

Once you've got the services up and running, navigate to [http://localhost/](http://localhost/) to see the frontend.

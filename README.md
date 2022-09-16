# pkgsite-live <!-- omit in toc -->

<h4 align="center"><b>A docker image of pkgsite to see local modules documentations with livereload</b></h4>

&nbsp;
<div align="center">

<a href="https://github.com/mauroalderete/pkgsite-local-live/blob/main/LICENSE">
	<img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg">
</a>
<a href="https://github.com/mauroalderete/pkgsite-local-live/blob/main/CODE_OF_CONDUCT.md">
	<img alt="Contributor covenant: 2.1" src="https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg">
</a>
<a href="https://semver.org/">
	<img alt="Semantic Versioning: 2.0.0" src="https://img.shields.io/badge/Semantic--Versioning-2.0.0-a05f79?logo=semantic-release&logoColor=f97ff0">
</a>
<a href="https://pkg.go.dev/github.com/mauroalderete/pkgsite-local-live">
	<img src="https://pkg.go.dev/badge/github.com/mauroalderete/pkgsite-local-live.svg" alt="Go Reference">
</a>

[![Tests](https://github.com/mauroalderete/pkgsite-local-live/actions/workflows/tests.yml/badge.svg)](https://github.com/mauroalderete/pkgsite-local-live/actions/workflows/tests.yml)

<a href="https://github.com/mauroalderete/pkgsite-local-live/issues/new/choose">Report Bug</a>
·
<a href="https://github.com/mauroalderete/pkgsite-local-live/issues/new/choose">Request Feature</a>

<a href="https://twitter.com/intent/tweet?text=👋%20Check%20this%20amazing%20repo%20https://github.com/mauroalderete/pkgsite-local-live,%20created%20by%20@_mauroalderete%0A%0A%23golang%20%23pkgsite%20%23docker%20✌️">
	<img src="https://img.shields.io/twitter/url?label=Share%20on%20Twitter&style=social&url=https%3A%2F%2Fgithub.com%2Fatapas%2Fmodel-repo">
</a>
</div>

&nbsp;

# Content <!-- omit in toc -->
- [:wave: Introducing `pkgsite-local-live`](#wave-introducing-pkgsite-local-live)
- [:clamp: Use](#clamp-use)
  - [Run](#run)
  - [Ports](#ports)
  - [Volumes](#volumes)
  - [Examples](#examples)
- [:rocket: Upcomming Features](#rocket-upcomming-features)
- [:hammer: How to Set up `pkgsite-local-live` for Development?](#hammer-how-to-set-up-pkgsite-local-live-for-development)
  - [Test reloader service](#test-reloader-service)
  - [Build image](#build-image)
- [:handshake: Contributing to `pkgsite-local-live`](#handshake-contributing-to-pkgsite-local-live)
- [:pray: Support](#pray-support)

&nbsp;
# :wave: Introducing `pkgsite-local-live`
`pkgsite-local-live` is a docker image that maintains a pkgsite instance up with all go modules stored in the folder `$GOPATH/src` of the container. A watcher looks at any change in the go files to know when to restart the pkgsite instance and back reload all open browser views.

Binding your local `$GOPATH/src` you can use `pkgsite-local-live` to query the documentation from the local projects stored in your personal workspace, at the same time you can view the changes that occur while you are working on them in real-time.

```mermaid
graph TB
    subgraph local[Local workspace]
        direction TB
        workspace[fa:fa-folder Workspace]
        developer[fa:fa-user Go developer]
        browser[fa:fa-eye Browser]

        developer --->|Edits source files| workspace
        browser --->|Shows last changes<br>on documentation| developer
    end

    subgraph container[Pkgsite Local Live]
        direction LR
        subgraph volumes[Volumes]
            modules[Go modules]
        end
        subgraph services[Services]
            watcher[Watcher service]
            reloader[Reloader server]
            pkgsite[Pkgsite server]
        end
    end

    workspace -->|Volume binding| modules
    modules -->|Listens for any change<br>on go files| watcher
    modules -->|Loads once at the start<br>the documentation<br>from golang source files| pkgsite
    watcher -->|Restarts the pkgsite instance<br>to reload the last changes<br>on documentation| pkgsite
    watcher -->|Emits a request<br>to handle a reload event| reloader
    pkgsite -->|Gets documentation pages| reloader
    reloader -->|Serves documentation pages<br>with a live-reload system injected| browser
```

In the next diagram you can look the three main sequences that it's contained in `pkgsite-local-live`.

The first sequence block shows how are executed the internal service. A server named Reloader Server provide a proxy to pkgsite sources at the same time handle the reload system event to implement a live-reload feature. A watcher process is charge of listen any change on workspace folder that can contains change on documentation pages.

The next sequence shows us the process involucred when a developer visit a documentation page with his browser.

The last sequence shows us the steps launched when a go file is modified, created or deleted.

```mermaid
sequenceDiagram
    autonumber
    rect rgb(191, 223, 255)
        note right of Entrypoint: Starts container
        Entrypoint-)ReloaderServer: Starts reloader server<br>with a proxy and interceptors
        Entrypoint-)Watcher: Starts watcher service<br>to listen any change on workspace folder
        Watcher->>PkgsiteServer: Starts pkgsite server instance<br>loading all go modules in workspace folder
        activate PkgsiteServer
        Workspace->>PkgsiteServer: Loads Go modules
        deactivate PkgsiteServer
    end

    rect rgb(191, 230, 223)
        note right of Developer: Developer queries documentation page
        actor Developer
        Developer->>Browser: Queries documentation about a go module
        activate Browser
        Browser->>ReloaderServer: Gets documentation pages
        activate ReloaderServer
        ReloaderServer->>PkgsiteServer: Requests documentation pages
        activate PkgsiteServer
        PkgsiteServer->>ReloaderServer: Responses with some resources
        deactivate PkgsiteServer
        ReloaderServer->>Browser: Responses with documentation pages with live-reload system
        deactivate ReloaderServer
        Browser-->>ReloaderServer: Connects to websocket to listen reload signal
        activate ReloaderServer
        ReloaderServer-->>Browser: Accepts websocket connection
        deactivate ReloaderServer
        Browser->>Developer: Shows documentation resource
        deactivate Browser
    end

    rect rgb(223, 230, 191)
        note right of Developer: Changes some go file on workspace
        Developer-)Workspace: Saves changes on a go file
        Workspace->>Watcher: Changes detected on go files
        activate Watcher
        Watcher-)PkgsiteServer: Restarts the server instance
        activate Watcher
        activate PkgsiteServer
        Workspace->>PkgsiteServer: Loads Go modules
        deactivate PkgsiteServer
        note right of Watcher: Waits 1 second
        Watcher-)ReloaderServer: Requests a reload the all clients
        deactivate Watcher
        deactivate Watcher
        activate ReloaderServer
        ReloaderServer-)Browser: Sends a websocket message to refresh the views
        deactivate ReloaderServer
        activate Browser

        note left of Browser: Reloads when the pages rendered in client recieved the reload websocket message

        Browser->>ReloaderServer: Gets documentation pages
        activate ReloaderServer
        ReloaderServer->>PkgsiteServer: Requests documentation pages
        activate PkgsiteServer
        PkgsiteServer->>ReloaderServer: Responses with some resources
        deactivate PkgsiteServer
        ReloaderServer->>Browser: Responses with documentation pages with live-reload system
        deactivate ReloaderServer
        Browser-->>ReloaderServer: Connects to websocket to listen reload signal
        activate ReloaderServer
        ReloaderServer-->>Browser: Accepts websocket connection
        deactivate ReloaderServer
        Browser->>Developer: Shows documentation resource
        deactivate Browser
    end
```

Please, look at [Contributing to `pkgsite-local-live`](#handshake-contributing-to-pkgsite-local-live) to choose the way to collaborate with you feel better.

# :clamp: Use

## Run

```bash
docker run -v $GOPATH/src:/go/src -p 8080:80 mauroalderete/pkgsite-local-live:latest
```

## Ports

Exposes the port 80 to access to pkgsite instance with all modules loaded.

## Volumes

`pkgsite-local-live` searches the modules in `/go/src` path. You must provide a source that will contains the go modules.

If the volume source doesn't have any go module, the pkgsite instance will end with an error and you cannot see anything through the port. This state will maintain this way to a 'go.mod' file will be found.

It is expected that the source volume contains many go modules, each one in its own folder.

```
./myvolume
  |- project-1
    |- go.mod
  |- project-2
    |- go.mod
  ...
  |- project-n
    |- go.mod
```

## Examples

```bash
docker run -v $GOPATH/src:/go/src -p 8080:80 mauroalderete/pkgsite-local-live:latest
```

Configures a container to load in pkgsite instance all modules stored in the golang standard workspace. Binds the port 8080 to access to pkgsite website.

# :rocket: Upcomming Features

`pkgsite-local-live` has all the potential to grow further. Here are some of the upcoming features planned (not in any order),

- ✔️ Filter modules to load. You will can filter the modules that you want to load by pkgsite instance using a yaml file configure.
- ✔️ Index. You will can enable a index in the home page to view all modules loaded and visit to speedly.

# :hammer: How to Set up `pkgsite-local-live` for Development?

You set up `pkgsite-local-live` locally with a few easy steps.

1. Clone the repository

```bash
git clone https://github.com/mauroalderete/pkgsite-local-live
cd pkgsite-local-live
```

## Test reloader service

```bash
got=go test -v ./... -coverprofile=coverage.out -covermode=count && go tool cover -html=coverage.out
```
## Build image

```bash
docker build -t <username>/pkgsite-local-live:<tag> .
```
# :handshake: Contributing to `pkgsite-local-live`

Any kind of positive contribution is welcome! Please help us to grow by contributing to the project.

If you wish to contribute, you can work on any [issue](https://github.com/mauroalderete/pkgsite-local-live/issues/new/choose) or create one on your own. After adding your code, please send us a Pull Request.

> Please read [`CONTRIBUTING`](CONTRIBUTING.md) for details on our [`CODE OF CONDUCT`](CODE_OF_CONDUCT.md), and the process for submitting pull requests to us.

# :pray: Support

We all need support and motivation. `pkgsite-local-live` is not an exception. Please give this project a :star: start to encourage and show that you liked it. Don't forget to leave a :star: star before you move away.

If you found the app helpful, consider supporting us with a coffee.

<div align="center">
<a href='https://cafecito.app/mauroalderete' rel='noopener' target='_blank'><img srcset='https://cdn.cafecito.app/imgs/buttons/button_6.png 1x, https://cdn.cafecito.app/imgs/buttons/button_6_2x.png 2x, https://cdn.cafecito.app/imgs/buttons/button_6_3.75x.png 3.75x' src='https://cdn.cafecito.app/imgs/buttons/button_6.png' alt='Invitame un café en cafecito.app' /></a>
</div>

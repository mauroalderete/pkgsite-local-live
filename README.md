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
  - [Build](#build)
- [:hamburger: Built With](#hamburger-built-with)
- [:handshake: Contributing to `pkgsite-local-live`](#handshake-contributing-to-pkgsite-local-live)
- [:pray: Support](#pray-support)

&nbsp;
# :wave: Introducing `pkgsite-local-live`
`pkgsite-local-live` is a docker image that maintains a pkgsite instance up with all go modules stored in a folder loaded. A watcher looks at any change in the go files to know when to restart the pkgsite instance and back reload your browser view.

With `pkgsite-local-live` you can query the documentation from the local projects stored in your personal workspace, at the same time that you can view the changes that occur while you working on them in real-time.

For the moment, this project is in an early development phase, so all versions are unstable.

You can take a look at the [upcoming features](#rocket-upcomming-features) to know more about the `pkgsite-local-live` future.

Hey! don't be discouraged, you can help me to carry out this project in many ways, contributing with new features, reporting bugs, sharing in your social networks or supporting with a :star:

Please, look at [Contributing to `pkgsite-local-live`](#handshake-contributing-to-pkgsite-local-live) to choose the way to collaborate with you feel better.

# :clamp: Use

## Run

```bash
docker run -v <GO WORKSPACE>:/go/src -p <PORT>:3000 mauroalderete/pkgsite-local-live:latest
```

## Ports

Exposes the port 3000 to access to pkgsite instance with all modules loaded.

## Volumes

`pkgsite-local-live` searches the modules in `/go/src` path. You must provide a source that will contain the go modules.

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
docker run -v $GOPATH/src:/go/src -p 8080:3000 mauroalderete/pkgsite-local-live:latest
```

Configures a container to load in pkgsite instance all modules stored in the golang standard workspace. Binds the port 8080 to access to pkgsite website.

# :rocket: Upcomming Features

`pkgsite-local-live` has all the potential to grow further. Here are some of the upcoming features planned (not in any order),

- ✔️ Liverelaod. Any change in your go files will show updating in your browser.
- ✔️ Auto-scan modules. Only you will need to configure a volume with your workspace and all go modules will be loaded.

# :hammer: How to Set up `pkgsite-local-live` for Development?

You set up `pkgsite-local-live` locally with a few easy steps.

1. Clone the repository

```bash
git clone https://github.com/mauroalderete/pkgsite-local-live
cd pkgsite-local-live
```

## Build

```bash
docker build -t <username>/pkgsite-local-live:<tag> .
```

# :hamburger: Built With

- [Docker](https://www.docker.com/) 20.10.17

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

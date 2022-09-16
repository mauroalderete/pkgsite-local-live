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
- [Introducing `pkgsite-local-live`](#introducing-pkgsite-local-live)
- [Use](#use)
  - [Run](#run)
  - [Ports](#ports)
  - [Volumes](#volumes)
  - [Examples](#examples)
- [Upcomming Features](#upcomming-features)
- [How to Set up `pkgsite-local-live` for Development?](#how-to-set-up-pkgsite-local-live-for-development)
  - [Test reloader service](#test-reloader-service)
  - [Build image](#build-image)
- [Contributing to `pkgsite-local-live`](#contributing-to-pkgsite-local-live)
- [Support](#support)

&nbsp;
# Introducing `pkgsite-local-live`
`pkgsite-local-live` is a docker image that maintains a pkgsite instance up with all go modules stored in the folder `$GOPATH/src` of the container. A watcher looks at any change in the go files to know when to restart the pkgsite instance and back reload all open browser views.

Binding your local `$GOPATH/src` you can use `pkgsite-local-live` to query the documentation from the local projects stored in your personal workspace, at the same time you can view the changes that occur while you are working on them in real-time.

[![](https://mermaid.ink/img/pako:eNqFVE1v2zAM_SuCzknROp8Ohh6KFb3kMCTDCizZQbFoW6steZKcLEvy30dZivM5TAdboh8f-UjKO5ooDnRCM82qnHx9WUqCy9QrbyhUworF1D3JRukPU7EEfniQW1xoSKxQsnV1q0UuUjZJWTdVBQdN3u8RwBoKVYEO0Nog8E2d7GfYlVYb0yJhC-TFWxBzh5F0u93n_SsX1hCjap0ASUUBZn_K74bb-8xzPJCCGUuSnMkMzKeVfkaRXCV1CdIyJ3l_CuV5QPJjIm0BE4VoITHrLx-ZERaIL-ZUrO_XcTo7WVuWtSowrll88-8zR7dKxWvUtcCyhe0ZoEnqhhClrkWCPvOwuaLcMJvkmPS7fx_xVyiN6hl2djELmwZ30TO3Kq-8rcANqM3xrIJtj4jriBdOVkJyIbP9UaeHhkMDnApjQRqSKk2Y3Ib-hfZlqp0Ar-uOP-owREmMyyyxOWZrmbaOwB0uBsAZU61K5C0wytWQBdVBTaijCzGDhtI07AFFhEQjRm0CqVDZBvH_MfxnpNfSzT5Dtl81Bg3cSMULaMxNEBxiafdtMz3NMS9H8wbIchGTVCxzGi99jqfGyQ0W3HVzaWyEzTGDAm9BN6Rhtti5EgvxE28C8P3xStIOLUGXTHD8Te1cpCXFupSwpBPcckhZXdglXcoDQllt1XwrE4q_iMJAh9YVZxY-C4aDX7bWikk62dHfdBINBg-9cdTvR3E8HMbxuN-hWzrpRmgexaN-PBiOH5_Gw0OH_lEKGZ4eHsfROO73okE0GPb6o1FD9735aHUNh79RrMnH)](https://mermaid.live/edit#pako:eNqFVE1v2zAM_SuCzknROp8Ohh6KFb3kMCTDCizZQbFoW6steZKcLEvy30dZivM5TAdboh8f-UjKO5ooDnRCM82qnHx9WUqCy9QrbyhUworF1D3JRukPU7EEfniQW1xoSKxQsnV1q0UuUjZJWTdVBQdN3u8RwBoKVYEO0Nog8E2d7GfYlVYb0yJhC-TFWxBzh5F0u93n_SsX1hCjap0ASUUBZn_K74bb-8xzPJCCGUuSnMkMzKeVfkaRXCV1CdIyJ3l_CuV5QPJjIm0BE4VoITHrLx-ZERaIL-ZUrO_XcTo7WVuWtSowrll88-8zR7dKxWvUtcCyhe0ZoEnqhhClrkWCPvOwuaLcMJvkmPS7fx_xVyiN6hl2djELmwZ30TO3Kq-8rcANqM3xrIJtj4jriBdOVkJyIbP9UaeHhkMDnApjQRqSKk2Y3Ib-hfZlqp0Ar-uOP-owREmMyyyxOWZrmbaOwB0uBsAZU61K5C0wytWQBdVBTaijCzGDhtI07AFFhEQjRm0CqVDZBvH_MfxnpNfSzT5Dtl81Bg3cSMULaMxNEBxiafdtMz3NMS9H8wbIchGTVCxzGi99jqfGyQ0W3HVzaWyEzTGDAm9BN6Rhtti5EgvxE28C8P3xStIOLUGXTHD8Te1cpCXFupSwpBPcckhZXdglXcoDQllt1XwrE4q_iMJAh9YVZxY-C4aDX7bWikk62dHfdBINBg-9cdTvR3E8HMbxuN-hWzrpRmgexaN-PBiOH5_Gw0OH_lEKGZ4eHsfROO73okE0GPb6o1FD9735aHUNh79RrMnH)

In the next diagram you can look the three main sequences that it's contained in `pkgsite-local-live`.

The first sequence block shows how are executed the internal service. A server named Reloader Server provide a proxy to pkgsite sources at the same time handle the reload system event to implement a live-reload feature. A watcher process is charge of listen any change on workspace folder that can contains change on documentation pages.

The next sequence shows us the process involucred when a developer visit a documentation page with his browser.

The last sequence shows us the steps launched when a go file is modified, created or deleted.

[![](https://mermaid.ink/img/pako:eNrtV01v2kAQ_Ssrn4oEUTAfwahCapMqlx7acIhUcVnWg72KvevuriE0yn_vrG1sgzcRyaVSFR8QzHrefL5n_OQxGYI39zT8zkEwuOE0UjRdCYIXzY0UeboGVf5WwAxR0frTMBj2ie-P8GMy6ZWH9hLSAFE8ig2RG_JNGLXPJBdmTpaGKqMJk8JQLg6A9mruGvTuIJE0BLUEtQVVe6nKTHRh_7xWix03MaEkU_JxT6gICfqDYpAZqbQb_J4aFrdQd-XvApQzsKhGkoRrAwIh94TFVERApCA7qR50RhmQjUzCdvYV6GCx-PEQaW7gJPWstFaZY5baUFEGszVxERGaJCSSJJVhnoDGW14JR5nhW4p4R9Fa6Rw8uwl9x3Ca3NaBGqcQXoMFEa6Ec_6jy2IJXpz_DWxxcJkNXn8luGaKY5mhZHkKuA2GY4MzGsFRlVI1Ps1BbcLyviq50xb7pxORrmVucEPqzjqaWEE0J5UB0U838RaMK2ftQD12bW44tncHdGcpqM8O88IGHJkdhdyBzqTQ2K-CQlqmODHQMkfynL0TznLqgZyEcJRTHiR8C4OS20TvkXapM_5L_TwMy1HktRQCd1UTJPQO1lqyBzCkYfchKI8ETd4_wnbRX5gVH90Kx8oksOh3lbVYtPizjNF40snD1JzoR6vtonCp3pbCSOZzKHxdqKEuNwZZteHJsTa6aNqrBQlroFt0ZxWMpegBxq1ftV4fIodgMHsIrW_l6SJG5dfV6F6HcLoUaRN3BPoc3H8gxt0J1V26pxxLGWIhuHiho_quElR6Qw98sH2wTyOWcNwyd1KdLrx29lZC9Wo6LXFnbWYNnVLQGrXDsljBBpc_LvLdctjpNxGsy5GT3iawKVrbErSkmNcuRu2wQUsRU5gjKLD_PaqWWXpxXP6ylVVXOyW0A348cT6eOP_nE8freymolPIQXy-erHnlIStSWHlz_BrChuaJWXkr8Yy32leN5V4wb76hiYa-l2cholavI7U1o8KbP3mP3hxfPC5GM3889oNgOg2C2bjv7b35wEfzVXA1DibT2eVwNn3ue3-kRIThxeXMnwXjkT_xJ9PR-OqqgPtVHBqVw_NfdEBkqQ)](https://mermaid.live/edit#pako:eNrtV01v2kAQ_Ssrn4oEUTAfwahCapMqlx7acIhUcVnWg72KvevuriE0yn_vrG1sgzcRyaVSFR8QzHrefL5n_OQxGYI39zT8zkEwuOE0UjRdCYIXzY0UeboGVf5WwAxR0frTMBj2ie-P8GMy6ZWH9hLSAFE8ig2RG_JNGLXPJBdmTpaGKqMJk8JQLg6A9mruGvTuIJE0BLUEtQVVe6nKTHRh_7xWix03MaEkU_JxT6gICfqDYpAZqbQb_J4aFrdQd-XvApQzsKhGkoRrAwIh94TFVERApCA7qR50RhmQjUzCdvYV6GCx-PEQaW7gJPWstFaZY5baUFEGszVxERGaJCSSJJVhnoDGW14JR5nhW4p4R9Fa6Rw8uwl9x3Ca3NaBGqcQXoMFEa6Ec_6jy2IJXpz_DWxxcJkNXn8luGaKY5mhZHkKuA2GY4MzGsFRlVI1Ps1BbcLyviq50xb7pxORrmVucEPqzjqaWEE0J5UB0U838RaMK2ftQD12bW44tncHdGcpqM8O88IGHJkdhdyBzqTQ2K-CQlqmODHQMkfynL0TznLqgZyEcJRTHiR8C4OS20TvkXapM_5L_TwMy1HktRQCd1UTJPQO1lqyBzCkYfchKI8ETd4_wnbRX5gVH90Kx8oksOh3lbVYtPizjNF40snD1JzoR6vtonCp3pbCSOZzKHxdqKEuNwZZteHJsTa6aNqrBQlroFt0ZxWMpegBxq1ftV4fIodgMHsIrW_l6SJG5dfV6F6HcLoUaRN3BPoc3H8gxt0J1V26pxxLGWIhuHiho_quElR6Qw98sH2wTyOWcNwyd1KdLrx29lZC9Wo6LXFnbWYNnVLQGrXDsljBBpc_LvLdctjpNxGsy5GT3iawKVrbErSkmNcuRu2wQUsRU5gjKLD_PaqWWXpxXP6ylVVXOyW0A348cT6eOP_nE8freymolPIQXy-erHnlIStSWHlz_BrChuaJWXkr8Yy32leN5V4wb76hiYa-l2cholavI7U1o8KbP3mP3hxfPC5GM3889oNgOg2C2bjv7b35wEfzVXA1DibT2eVwNn3ue3-kRIThxeXMnwXjkT_xJ9PR-OqqgPtVHBqVw_NfdEBkqQ)

Please, look at [Contributing to `pkgsite-local-live`](#handshake-contributing-to-pkgsite-local-live) to choose the way to collaborate with you feel better.

# Use

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

# Upcomming Features

`pkgsite-local-live` has all the potential to grow further. Here are some of the upcoming features planned (not in any order),

- ✔️ Filter modules to load. You will can filter the modules that you want to load by pkgsite instance using a yaml file configure.
- ✔️ Index. You will can enable a index in the home page to view all modules loaded and visit to speedly.

# How to Set up `pkgsite-local-live` for Development?

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
# Contributing to `pkgsite-local-live`

Any kind of positive contribution is welcome! Please help us to grow by contributing to the project.

If you wish to contribute, you can work on any [issue](https://github.com/mauroalderete/pkgsite-local-live/issues/new/choose) or create one on your own. After adding your code, please send us a Pull Request.

> Please read [`CONTRIBUTING`](CONTRIBUTING.md) for details on our [`CODE OF CONDUCT`](CODE_OF_CONDUCT.md), and the process for submitting pull requests to us.

# Support

We all need support and motivation. `pkgsite-local-live` is not an exception. Please give this project a start to encourage and show that you liked it.

If you found the app helpful, consider supporting us with a coffee.

<div align="center">
<a href='https://cafecito.app/mauroalderete' rel='noopener' target='_blank'><img srcset='https://cdn.cafecito.app/imgs/buttons/button_6.png 1x, https://cdn.cafecito.app/imgs/buttons/button_6_2x.png 2x, https://cdn.cafecito.app/imgs/buttons/button_6_3.75x.png 3.75x' src='https://cdn.cafecito.app/imgs/buttons/button_6.png' alt='Invitame un café en cafecito.app' /></a>
</div>

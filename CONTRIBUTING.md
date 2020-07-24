# Contributing

By participating to this project, you agree to abide our [code of
conduct](/CODE_OF_CONDUCT.md).

## Setup your machine

`hub` is written in [Go](https://golang.org/).

Prerequisites are:

- Build:
  - `make`
  - [Go 1.9+](http://golang.org/doc/install)

Clone `hub` from source into some path:

```sh
git clone git@github.com:leandro-lugaresi/hub.git
cd hub
```

If you created a fork clone your fork and add my repository as a upstream remote:

```sh
git clone git@github.com:{your-name}/hub.git
cd hub
git remote add upstream git@github.com:leandro-lugaresi/hub.git
```

### Install

Install the build and lint dependencies:

```sh
$ make setup
```

A good way of making sure everything is all right is running the test suite:

```sh
$ make test
```

## Test your change

When you are satisfied with the changes, we suggest you run:

```sh
$ make ci
```

Which runs all the linters and tests.

## Submit a pull request

Push your branch to your `example` fork and open a pull request against the
main branch.

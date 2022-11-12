# The Go SSA Playground

[![go-recipes](https://raw.githubusercontent.com/nikolaydubina/go-recipes/main/badge.svg?raw=true)](https://github.com/nikolaydubina/go-recipes)

https://golang.design/gossa

A tool for exploring Go's SSA intermediate representation.

![](./public/assets/screen.png)

## Deployment

There are two approaches to use the SSA Playground: native execution
or Docker-based deployment.

To execute natively, you can just use:

```bash
$ make
```

If you have Docker, then you can use:

```bash
$ make build # build the docker image
$ make up    # run/update for latest image
```

Then access http://localhost:6789/gossa.

## Contribution

We would love to hear your feedback. Submit [an issue](https://github.com/golang-design/ssaplayground/issues/new) to tell bugs or feature requests.

## License

GNU GPLv3 &copy; The [golang.design](https://golang.design) Authors. Originally created by [Changkun Ou](https://changkun.de).

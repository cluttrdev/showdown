# Showdown

Showdown is a live Markdown previewer.

It uses [fsnotify][github-fsnotify] to watch a given file for changes,
renders the content to HTML using [markdown][github-gomarkdown]
and serves the output on your local machine.

For an update on the latest changes please have a look at the
[CHANGELOG](./CHANGELOG.md).

## Installation

To install `showdown` you can download a prebuild binary from the 
[releases](https://github.com/cluttrdev/showdown/releases) page.

E.g. if you're on Linux:

```shell
# determine latest release
RELEASE_TAG=$(curl -sSL https://api.github.com/repos/cluttrdev/showdown/releases/latest | jq -r '.tag_name')
# download release archive
curl -sSL -O https://github.com/cluttrdev/showdown/releases/download/${RELEASE_TAG}/showdown_${RELEASE_TAG}_linux_x86_64.tar.gz
# extract binary
tar -zxf showdown_${RELEASE_TAG}_linux_x86_64.tar.gz showdown
# install it
install ./showdown /usr/local/bin/showdown
```

Alternatively, you can install it using the standard Go tools.

```shell
go install github.com/cluttrdev/showdown@latest
```

## Usage

To preview a Markdown formatted file `example.md` simply run

```shell
showdown example.md
```

This will render the file content as HTML and serve it under <http://localhost:1337/>

Run `showdown --help` for more information.

## Acknowledgement

This project was inspired by [livedown][github-livedown].

## License

This project is released under the [MIT License](./LICENSE)

[github-fsnotify]: https://github.com/fsnotify/fsnotify
[github-gomarkdown]: https://github.com/gomarkdown/markdown
[github-livedown]: https://github.com/shime/livedown

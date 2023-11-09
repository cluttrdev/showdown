# Showdown

Showdown is a live Markdown previewer.

It uses [fsnotify][github-fsnotify] to watch a given file for changes,
renders the content to HTML using [markdown][github-gomarkdown]
and serves the output on your local machine.

For an update on the latest changes please have a look at the
[CHANGELOG](./CHANGELOG.md).

## Installation

To install `showdown` you can download a 
[prebuilt binary][prebuilt-binaries] that matches your system, e.g.

```shell
OS=linux
ARCH=amd64
# Download
RELEASE_TAG=$(curl -sSfL https://api.github.com/repos/cluttrdev/showdown/releases/latest | jq -r '.tag_name')
curl -sSfL https://github.com/cluttrdev/showdown/releases/download/${RELEASE_TAG}/showdown_${RELEASE_TAG}_${OS}_${ARCH}.tar.gz -o /tmp/showdown.tar.gz
# Install
tar -xozf /tmp/showdown.tar.gz showdown
install ./showdown ~/.local/bin/showdown
```

Alternatively, if you have the [Go][go-install] tools installed on your
machine, you can use

```shell
go install github.com/cluttrdev/showdown@latest
```

## Usage

To preview a Markdown formatted file, e.g. this project's `README.md`, simply
run

```shell
showdown README.md
```

This will render the file content as HTML, serve it on <http://127.0.0.1:1337/>
and update the preview on changes.

Run `showdown --help` for more information.

## Acknowledgement

This project was inspired by [livedown][github-livedown].

## License

This project is released under the [MIT License](./LICENSE)

## Troubleshooting

#### No live update when editing in (neo)vim

If you're previewing a file while editing it using (neo)vim it might not get
updated on writes. This is probably due to how backups are configured in your
setup (see `:help backup`).

> If you write to an existing file (but do not append) while the 'backup',
> 'writebackup' or 'patchmode' option is on, a backup of the original file is
> made. The file is either copied or renamed (see 'backupcopy').

If `backupcopy` is set to `"no"` (or `"auto"`) the original file is renamed and
a new one is written to. Thus, the backup file is checked for changes instead
of the newly created file.

To fix this, consider setting the `backupcopy` option to `"yes"`. This will
make a copy of the file and overwrite the original one.

<!-- Links -->
[github-fsnotify]: https://github.com/fsnotify/fsnotify
[github-gomarkdown]: https://github.com/gomarkdown/markdown
[github-livedown]: https://github.com/shime/livedown
[go-install]: https://go.dev/doc/install
[prebuilt-binaries]: https://github.com/cluttrdev/showdown/releases/latest

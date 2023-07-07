This project is a fork from https://github.com/fatih/hclfmt


# hclfmt

hclfmt is a command to format and prettify HCL files. It's similar to the
popular `gofmt` command. Hook it with your favourite editor or use it from the
command line.

## Install

If you have Go installed just do:

```bash
go install github.com/magodo/hclfmt@master
```

## Editor integration

* [vim-hclfmt plugin](https://github.com/fatih/vim-hclfmt)
* [atom-hclfmt](https://atom.io/packages/hclfmt)

## Usage

The usage is similar to `gofmt`. If you pass a file it prints the formatted
output to std output:

```bash
$ hclfmt config.hcl
```

You can pass the `-w` flag to directly overwrite your file:

```bash
$ hclfmt -w config.hcl
```

If no arguments are passed, it excepts the input from standard input.


## License

The BSD 3-Clause License - see
[`LICENSE`](https://github.com/fatih/hclfmt/blob/master/LICENSE) for more
details


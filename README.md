## What
Numbers any valid file path in stdin (one per line). Copies chosen file names to the clipboard.

## Why
For example you want to undo some changes you made to a git project:
```sh
$ git st
## master
 M README.md
 M somwhere/in/a/deeply/nested/directory/program.go
```

You can't be bothered reaching for the mouse to copy this file name so you pipe the output to numberwang:

```sh
$ git st | nw
## master
{1}  M {README.md}
{2}  M {somwhere/in/a/deeply/nested/directory/program.go}

to clipboard: 
```
You now got prompted to choose file names to copy to the clipboard. You choose "2".

```sh
$ git st | nw
## master
{1}  M {README.md}
{2}  M {somwhere/in/a/deeply/nested/directory/program.go}

to clipboard: 2
nw: wrote "somwhere/in/a/deeply/nested/directory/program.go " to clipboard
```

Now you can simply paste the file name(s) you selected when performing a checkout:

```sh
$ git checkout somwhere/in/a/deeply/nested/directory/program.go 
```

## Usage with git
I recommend a git alias that preserves colored output, for example:

```st = -c color.status=always status -sb```

Or if you want to go all-in and always call numberwang:

```snw = ! git -c color.status=always status -sb | nw```

Other commands might have similar options to preserve color.

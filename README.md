Numbers any valid file path in stdin (one per line). Copies chosen file names to the clipboard.

```
$ git st
## master
 M README.md
 M program.go

$ git st | nw
## master
{1}  M {README.md}
{2}  M {program.go}

to clipboard: 1 2
nw: wrote "README.md program.go " to clipboard

$ git checkout README.md program.go 
```

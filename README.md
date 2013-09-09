Numbers any valid file path in stdin (one per line). Copies chosen file names to the clipboard.

```
$ ls
nw*  nw.go

$ ls | nw
1 nw
2 nw.go

$ ls | nw 1 2
1 nw
2 nw.go

# the string 'nw nw.go ' is now on the clipboard
```

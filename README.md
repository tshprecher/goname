# goname

Goname properly formats the names of package level variables and constants according to golang conventions. For example, a variable named "GLOBAL_VARIABLE" will be renamed to "GlobalVariable". This does not modify access levels. Private variables remain private and public variables remain public.

Renaming is implemented via calls to `gorename`. To install gorename, run 
```
$ go install golang.org/x/tools/cmd/gorename
```

Make sure gorename exists in your $PATH.

## examples:

To format the names in a given package, run
```
$ goname path/to/golang/package
```

To identify malformed names without renaming, run
```
$ goname -l path/to/golang/package
```

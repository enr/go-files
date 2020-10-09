Go files
========

![CI Linux](https://github.com/enr/go-files/workflows/CI%20Nix/badge.svg)
![CI Windows](https://github.com/enr/go-files/workflows/CI%20Windows/badge.svg)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/enr/go-files)](https://pkg.go.dev/github.com/enr/go-files)
[![Go Report Card](https://goreportcard.com/badge/github.com/enr/go-files)](https://goreportcard.com/report/github.com/enr/go-files)

Shortcuts for operations pertaining files.

Import the library:

```Go
import (
    "github.com/enr/go-files/files"
)
```

Check if path exists (is a file or a directory):

NOTE: There are some good reasons why this function is not implemented in the standard lib (see https://github.com/golang/go/issues/1312)
but, if you have to simply check if a file exists in this very moment, it is useful to have a one liner.


```Go
e := files.Exists(fpath)
```

Check if path is a regular file:

```Go
e := files.IsRegular(fpath)
```

Check if path is a directory:

```Go
is := files.IsDir(fpath)
```

Check if path is a symlink:

```Go
e := files.IsSymlink(fpath)
```

Get the Sha1 checksum:

```Go
sha1sum, err := files.Sha1Sum(fpath)
if err != nil {
    // ...
}
```

Copy file to file:

```Go
err := files.Copy(srcpath, destpath)
if err != nil {
    // ...
}
```

Copy file in directory:

```Go
err := files.Copy(fpath, destpath)
if err != nil {
    // ...
}
```

Copy directory recursively:

```Go
err := files.CopyDir(source, destination)
if err != nil {
    // ...
}
```

Read file lines:

```Go
lines, err := files.ReadLines(fpath)
if err != nil {
    // ...
}
```

Process file lines:

```Go
files.EachLine(path, func(line string) error {
    fmt.Println(" # " + line)
    return nil
})
```


## License

Apache 2.0 - see LICENSE file.

Copyright 2014 go-files contributors

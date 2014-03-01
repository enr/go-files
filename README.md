Go files
========

[![Build Status](https://travis-ci.org/enr/go-files.png?branch=master)](https://travis-ci.org/enr/go-files)

Go library to easily manipulate files.

Import the library:

```Go
    import (
        "github.com/enr/go-files/files"
    )
```

Check if path exists (is a file or a directory):

```Go
    e := files.Exists(data.path)
```

Check if path is a regular file:

```Go
    e := files.IsRegular(data.path)
```

Check if path is a directory:

```Go
    is := files.IsDir(data.path)
```

Get the Sha1 checksum:

```Go
    sha1sum, err := files.Sha1Sum(data.path)
    if err != nil {
        // ...
    }
```

Copy file to file:

```Go
    err := files.Copy(data.path, of)
    if err != nil {
        // ...
    }
```

Copy file in directory:

```Go
    err := files.Copy(data.path, of)
    if err != nil {
        // ...
    }
```

License
-------

Apache 2.0 - see LICENSE file.

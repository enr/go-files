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

   Copyright 2014 go-bintray contributors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

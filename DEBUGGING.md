## Debugging Tips

**Golang debugging caveat**

It is straightforward to run a Golang debugger on the `lumogon` program (e.g.,
by directing the debugger to start at the `main` function in `lumogon.go` and
proceeding from there). Note, however, that once the `lumogon` program hands off
execution to code running in a container in the Docker infrastructure under
examination the debugger does not follow the execution inside the container.

**Debugging with Visual Studio Code**

If you use [Visual Studio Code](https://code.visualstudio.com/) as your
development environment, here is a sample `launch.json` file to use when running
`lumogon.go` under the debugger:

``` json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "remotePath": "",
            "port": 2345,
            "host": "127.0.0.1",
            "program": "${fileDirname}",
            "env": {},
            "args": ["container", "--all", "--debug", "--keep-harvesters"],
            "showLog": true
        }
    ]
}
```

**Dumping Golang data structures**

The [`go-spew`](https://github.com/davecgh/go-spew) Golang library implements a
pretty-printer that can recursively and verbosely describe the contents of Go
data structures.  Here is an example of a helper method which wraps `spew` and
provides tagged log output (including the file+line location of the calling
function) for dumping an arbitrary data structure:

`glide.yaml`:
``` yaml
package: github.com/davecgh/go-spew/spew
```


`debug.go`:
``` go
import "github.com/davecgh/go-spew/spew"

// LogDump outputs a delimited log section with the message string and a structured dump of the provided data object
func LogDump(message string, data interface{}) {
	logging.Stderr("\n\n====BEGIN======================================")
	_, file, line, _ := runtime.Caller(1)
	logging.Stderr(fmt.Sprintf("Debug - %s @ %s:%d:\n", message, file, line), spew.Sdump(data))
	logging.Stderr("====END=======================================\n\n")
}
```

(client code) `dpkg.go`:

``` go
import "github.com/puppetlabs/transparent-containers/cli/debug"

// ...

	for scanner.Scan() {
		txt := strings.Replace(scanner.Text(), "'", "", 2)
		if txt != "" {
    		debug.LogDump("data", data)
			data = append(data, txt)
		}
	}
```

Sample output:

```
====BEGIN======================================
[lumogon] 2017/04/14 22:24:40.573682 Debug - data @ /Users/rick/go/src/github.com/puppetlabs/transparent-containers/cli/capability/dpkg.go:49:
%!(EXTRA string=([]string) (len=8 cap=8) {
 (string) (len=12) "acl,2.2.52-2",
 (string) (len=18) "adduser,3.113+nmu3",
 (string) (len=13) "apt,1.0.9.8.4",
 (string) (len=19) "base-files,8+deb8u7",
 (string) (len=18) "base-passwd,3.5.37",
 (string) (len=18) "bash,4.3-11+deb8u1",
 (string) (len=19) "bsdutils,1:2.25.2-6",
 (string) (len=16) "bzip2,1.0.6-7+b3"
}
)
[lumogon] 2017/04/14 22:24:40.573698 ====END=======================================
```

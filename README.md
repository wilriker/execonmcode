# Description
This is a small extension to [DuetSoftwareFramework](https://github.com/christhamm/DuetSoftwareFramework)
to execute arbitrary system commands when a user-defined `M-Code` is encountered.

An example usage would be to execute system shutdown on the SBC when a e.g. `M7722` is run.

# Usage
```
$ execonmcode -help
Usage of ./execonmcode:
  -command value
        Command to execute
  -debug
        Print debug output
  -mCode value
        Code that will initiate execution of the command
  -socketPath string
        Path to socket (default "/var/run/dsf/dcs.sock")
  -trace
        Print underlying requests/responses
```

Starting from version 3 it is possible to provide an arbitrary number of `-mCode` + `-command` tuples. This way a
single instance of `execonmcode` can handle multiple commands. Side-effect is that there is no more default for `-mCode`.

## Parameters
`execonmcode` does provide a simple mechanism for parameter substitution. It is possible to pass string parameters to the
selected `M-Code` and have them inserted in the `-command`. In the command string they have to be single letters prefixed by
the percent-sign (`%`) and they must not be `G`, `M` or `T`.

All parameters that do not have a corresponding value in the `M-Code` will be forwarded as given.

### Parameters in systemd units
Since `%` is used to access systemd-specific variables in unit files it is
necessary to escape them by using double-percent, i.e. `%%`.

### Example
Run `execonmcode` as
```
$ execonmcode -command "mycommand %F %N %D"
```
Then you can use the following `M-Code` syntax to replace these parameters
```
M7722 F"my first parameter" N"my second parameter"
```
this will lead to an execution of
```
mycommand "my first parameter" "my second parameter" %D
```
Note that `%D` was passed as is since it was not given in the `M-Code`.

# Installation
* Download
* Rename to just `execonmcode`
* Make executable via `chmod +x execonmcode`
* Put it into any path of your `$PATH` e.g. `/usr/local/bin`
* Run it as `root` or with `sudo`
* Optional: use the `shutdownsbc.service` systemd unit (included in the repo) to run it at startup and let it shutdown the SBC (customize to your liking)

# Contribution
I am happy about comments, suggestions, bug reports, pull requests, etc. either here or in [the forum](https://forum.duet3d.com/topic/13194).

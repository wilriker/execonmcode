# Description
This is a small extension to [DuetSoftwareFramework](https://github.com/christhamm/DuetSoftwareFramework)
to execute arbitrary system commands when a user-defined `M-Code` is encountered.

An example usage would be to execute system shutdown on the SBC when a e.g. `M7722` is run.

# Usage
```
$ ./execonmcode --help
Usage of ./execonmcode:
  -command string
        Command to execute
  -mCode int
        Code that will initiate execution of the command (default 7722)
  -socketPath string
        Path to socket (default "/var/run/duet.sock")
```

# Installation
* Download
* Rename to just `execonmcode`
* Put it into `/usr/local/bin` (or any other path in your $PATH)
* Run it as `root`
* Optional: use the `shutdownsbc.service` systemd unit (included in the repo) to run it at startup and let it shutdown the SBC (customize to your liking)

# Contribution
I am happy about comments, suggestions, bug reports, pull requests, etc. either here or in [the forum](https://forum.duet3d.com/topic/13194).

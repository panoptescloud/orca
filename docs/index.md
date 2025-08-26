![coverage](https://raw.githubusercontent.com/panoptescloud/orca/badges/.badges/main/coverage.svg)

# Home

Whalecome!

## Initial installation

Head over to the [latest github release](https://github.com/panoptescloud/orca/releases/latest), and find the URL for the archive that best suits your system. Put into the script below, and you're good to go!

```sh
$ curl -o /tmp/orca.tar.gz -L "{url}"
$ (cd /tmp && tar -xzvf /tmp/orca.tar.gz)
$ chmod +x /tmp/orca
$ sudo mv /tmp/orca /usr/local/bin/orca
```

Optionally, clean up the files from `/tmp`...

```
$ rm /tmp/orca.tar.gz /tmp/LICENSE /tmp/README.md
```

### Updating

!!! tip "> 0.5.0"
    If you're using a version of at least 0.5.0, you can simply use the [self-update command](./CLI//orca_util_self-update.md). Run this to replace the currently installed binary with the latest version from github.

If you're on a version below this, then follow the same instructions as above, but `rm /usr/local/bin/orca`, first.
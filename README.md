# The Krypt Password Manager

![logo](https://user-images.githubusercontent.com/68542775/158347672-652d6682-2172-4332-ae75-399afdb670e4.png)

### Installation

``` bash
git clone https://github.com/raklaptudirm/krypt.git
cd krypt
make build
./bin/krypt # put this executable in your bin
```

### Commands

```bash
krypt add            # add a new password to krypt
krypt edit [regexp]  # edit a password whose name matches [regex]
krypt help [command] # get help about any provided krypt command
krypt list [regexp]  # list the passwords whose names matches [regex]
krypt login          # login to krypt with your master password
krypt logout         # logout of krypt
krypt rm [regexp]    # remove the password whose name matches [regex]
krypt version        # get the version information of krypt
```

### License

Krypt is licensed under the [Apache License, Version 2.0](https://opensource.org/licenses/Apache-2.0).

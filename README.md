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
krypt add            # add a password to your krypt database.
krypt edit [regexp]  # Edit a password whose name matches the provided regular expression.
krypt help [command] # Get help about the provided krypt command.
krypt list [regexp]  # List the passwords whose names match the provided regular expression.
krypt login          # Login to your krypt database with your master password.
krypt logout         # Logout of your krypt database.
krypt rm [regexp]    # Remove the password whose name matches the provided regular expression.
krypt version        # Get the version information of your krypt executable.
```
### License
Krypt is licensed under the [Apache License 2.0](https://github.com/raklaptudirm/krypt/blob/master/LICENSE).

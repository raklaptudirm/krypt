# v2.5.0

- Added first version of `CHANGELOG.md` to the repository.

# v2.4.0

- Implemented required measures and updates to bring `PBKDF-2` to action. Key generation and checksum now with `PBKDF-2`.

# v2.3.0

- Replaced `SHA-256` with `PBKDF-2`, using node.js `crypto` implementation: `crypto.pbkdf2Sync()` (wrapped by `PBKDF2_HASH()`).

# v2.2.0

- Fixed cryptographically insecure `Math.random()` with secure `crypto.randomInt()` (wrapped by `crypt.random()`).

# v2.1.0

- Faster database loading with removal of unnecessary reloads.

# v2.0.0

- `Wordy Passwords` introduced

# v1.3.0

- `2FA` introduced.

# v1.2.0

- `search` command to search for a password

# v1.1.0

- `list` command to list all passwords

# v1.0.0
- AES Encryption for password storage

- SHA Password checksum

- Password managing

- Password strength checker

- Password leak checking with [have i been pwned](https://haveibeenpwned.com)

- `unsafe` Strong password generator. **[ Uses cryptographically insecure `Math.random()` ]** [ Fixed in `v2.2.0` ]

- Password health checking.











​    

​    

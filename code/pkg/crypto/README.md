# Slog-Crypto

The encryption is done in two steps, first the given string from a log is encrypted with a password (randomly generated on each packet), then the password SHA-256 is encrypted again with the public RSA key and sent over the wire.

This process is done at the `slog-client` bin, when data arrives to the `slog-server`, the reverse unwrapping happens. The same process, but with a different key/pass is used to store data to the encrypted file (+ some other metadata).

+---------------+  +------------+
|  SHA256+AES   |  |    RSA     |
|  +---------+  |  |  +------+  |
|  |Plaintext|  |  |  | Pass |  |
|  +---------+  |  |  +------+  |
|  SHA256+AES   |  |    RSA     |
+---------------+, +------------+

The message sent is a string as this:

```bash
$AES256blob,$RSAencryptedPass
```

## Password Encryption

As the password might not be long enough or too long, a SHA256 hash of it is used to encrypt the plaintext, the encryption algorithm is AES-256.

A high level overview of the process looks like:

- Generate random password
- Create SHA256 hash of the password
- Use the hash as key for AES
- Encrypt plaintext with AES

## Key encryption

A high level overview of the process looks like:

- Read key (private/public)
- Use key to encrypt the SHA-256 of the password.

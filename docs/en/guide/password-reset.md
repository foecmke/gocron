# Password Reset (CLI)

When an administrator forgets their password — or gets locked out because their
two-factor authentication (2FA) device is lost — you can reset the password
directly from the command line using the built-in `resetpwd` command.

## Overview

The `resetpwd` command updates a user's password directly in the database.
It only initializes the configuration and database connection — it does **not**
start the web server, scheduler, or leader election — so it is safe to run while
gocron is running.

::: tip
`resetpwd` locates its configuration (`conf/app.ini`) relative to the gocron
**binary**, not the current directory, and automatically anchors a relative
SQLite database path to gocron's base directory. So you can run it from any
working directory — just use the **installed binary** (the same one the server
runs). For SQLite, it prints the database file it connects to; verify that line
matches the server's database before confirming.
:::

## Usage

```bash
gocron resetpwd [username]
```

| Argument / Flag | Description |
| --- | --- |
| `[username]` | The user to reset. Optional — defaults to `admin`. Matched against the username. |
| `--disable-2fa` | Also turn off two-factor authentication for that user. |

## Interactive flow

1. A confirmation prompt is shown (`y/N`). Anything other than `y`/`yes` cancels.
2. You are asked to enter the new password. Input is **hidden** (not echoed to the
   terminal).
3. If you type a password, you must enter it a second time to confirm. If the two
   entries differ, the operation is aborted and nothing is changed.
4. If you leave the password **blank** and press Enter, a random 12-character
   password is generated and printed to the screen.

## Examples

Reset the default `admin` account:

```bash
./gocron resetpwd
```

Reset a specific user:

```bash
./gocron resetpwd alice
```

Reset the password **and** disable 2FA (use this when an administrator is locked
out by a lost authenticator):

```bash
./gocron resetpwd admin --disable-2fa
```

### Inside Docker

If gocron runs in a container, exec into it and run the command against the
in-container binary:

```bash
docker exec -it gocron ./gocron resetpwd admin
```

To also clear 2FA:

```bash
docker exec -it gocron ./gocron resetpwd admin --disable-2fa
```

## Troubleshooting

- **`gocron is not installed; cannot reset password`** — the binary cannot find a
  completed installation (`conf/app.ini` + install lock) next to it. Use the
  installed gocron binary — the same one the server runs (in Docker, exec into the
  running container).
- **`user [xxx] not found`** — the username does not exist, or you are connected to
  a different database than expected. Double-check the username and the database
  configured in `conf/app.ini`.

::: warning
For security, prefer letting the tool generate a random password (leave the input
blank). After regaining access, log in and change the password through the web
interface, and re-enable 2FA if it was disabled.
:::

## Related Documentation

- [Security - Two-Factor Authentication (2FA)](./security-2fa)
- [Configuration](./configuration)

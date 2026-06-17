# Configuration

::: tip Configuration File Path
Configuration file path: `.gocron/conf/app.ini`

This configuration file is automatically created after system installation.
:::

## Configuration Options

### Basic Configuration

- **allow_ips**  
  Allowed client IPs, multiple IPs separated by commas, default is empty (no restriction)

- **app.name**  
  Application name

### Database Configuration

- **db.engine**  
  Database engine, supports `mysql` and `postgres`

- **db.host**  
  Database hostname

- **db.port**  
  Database port

- **db.user**  
  Database username

- **db.password**  
  Database password

- **db.charset**  
  Database charset

- **db.prefix**  
  Table prefix

- **db.database**  
  Database name

### API Configuration

- **api.key**  
  API interface key, API cannot be used without configuration

- **api.secret**  
  API interface secret, API cannot be used without configuration

- **api.sign.enable**  
  Enable signature verification

### TLS Configuration

- **enable_tls**  
  Enable TLS

- **ca_file**  
  CA certificate file path

- **cert_file**  
  Client certificate path

- **key_file**  
  Client private key path

## Configuration Example

```ini
[app]
name = gocron

[server]
allow_ips = 

[db]
engine = mysql
host = 127.0.0.1
port = 3306
user = root
password = 
database = gocron
charset = utf8mb4
prefix = 

[api]
key = your_api_key
secret = your_api_secret
sign.enable = true

[tls]
enable_tls = false
ca_file = 
cert_file = 
key_file = 
```

## Related Documentation

- [Security - TLS Mutual Authentication](./security-tls)
- [API Documentation](./api)

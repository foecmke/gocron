# Security - TLS Mutual Authentication

TLS mutual authentication provides a higher level of security for gocron, ensuring secure communication between clients and servers.

## What is TLS Mutual Authentication

TLS Mutual Authentication is a security mechanism that requires both client and server to provide certificates for authentication:
- **Server Authentication**: Client verifies the server's identity
- **Client Authentication**: Server verifies the client's identity

## Configuring TLS

### 1. Generate Certificates

First, you need to generate CA certificate, server certificate, and client certificate.

**Generate CA Certificate**:
```bash
# Generate CA private key
openssl genrsa -out ca.key 2048

# Generate CA certificate
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt
```

**Generate Server Certificate**:
```bash
# Generate server private key
openssl genrsa -out server.key 2048

# Generate server certificate signing request
openssl req -new -key server.key -out server.csr

# Sign server certificate with CA
openssl x509 -req -days 3650 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt
```

**Generate Client Certificate**:
```bash
# Generate client private key
openssl genrsa -out client.key 2048

# Generate client certificate signing request
openssl req -new -key client.key -out client.csr

# Sign client certificate with CA
openssl x509 -req -days 3650 -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt
```

### 2. Configure gocron

Add the following configuration to the configuration file `.gocron/conf/app.ini`:

```ini
[tls]
enable_tls = true
ca_file = /path/to/ca.crt
cert_file = /path/to/server.crt
key_file = /path/to/server.key
```

### 3. Restart Service

After modifying the configuration, restart the gocron service for the changes to take effect.

## Client Configuration

When using TLS mutual authentication, clients also need to configure certificates:

```bash
curl --cacert ca.crt --cert client.crt --key client.key https://gocron-server:5920
```

## Verify Configuration

You can use the following command to verify that the TLS configuration is correct:

```bash
openssl s_client -connect localhost:5920 -CAfile ca.crt -cert client.crt -key client.key
```

## Troubleshooting

### Common Issues

**Q: Cannot access after enabling TLS**
- Check if certificate paths are correct
- Confirm certificate file permissions
- Check gocron logs for detailed error information

**Q: Certificate verification failed**
- Confirm that CA certificate, server certificate, and client certificate are issued by the same CA
- Check if certificates are expired
- Verify that the Common Name (CN) in certificates is correct

## Best Practices

- **Regular certificate updates**: Certificates should be updated regularly to avoid expiration
- **Secure private key storage**: Private key files should have appropriate permissions to prevent leakage
- **Use strong encryption**: Use 2048-bit or higher RSA keys
- **Monitor certificate validity**: Set reminders to update certificates before expiration

## Related Documentation

- [Security - Two-Factor Authentication (2FA)](./security-2fa)
- [Configuration](./configuration)

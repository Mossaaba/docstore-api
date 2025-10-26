#!/bin/bash

# Generate self-signed SSL certificates for development/testing
# DO NOT use these certificates in actual production!

echo "Generating self-signed SSL certificates for development..."

# Generate private key
openssl genrsa -out key.pem 2048

# Generate certificate signing request
openssl req -new -key key.pem -out cert.csr -subj "/C=US/ST=State/L=City/O=Organization/OU=OrgUnit/CN=localhost"

# Generate self-signed certificate
openssl x509 -req -days 365 -in cert.csr -signkey key.pem -out cert.pem

# Clean up CSR file
rm cert.csr

echo "SSL certificates generated:"
echo "- cert.pem (certificate)"
echo "- key.pem (private key)"
echo ""
echo "⚠️  WARNING: These are self-signed certificates for development only!"
echo "   For production, use proper SSL certificates from a trusted CA."

# Local SSL setup

```
# Generate self-signed certificate: 
bash scripts/gen-cert.sh

# Set following line in /etc/hosts file:
127.0.0.1     redmine.local

# Run docker compose
docker compose up

# login to server by providing certificate
go run main.go login inline -s https://redmine.local -c certs/server.crt -u admin -p admin
```

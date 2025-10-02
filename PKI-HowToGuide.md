# 🔐 PKI Lab – Create Your Own CA and SSL Certificates

This guide walks through setting up your own **Root Certificate Authority (CA)** and issuing a **server certificate** for local development.  
By the end, you’ll be able to run a Go HTTPS server trusted by your browser (no warnings).

---

## 📂 File Layout

```
rootCA.key   # Root CA private key (keep secret, never commit)
rootCA.pem   # Root CA certificate (install in trust store)
root.cnf     # Root CA config

server.key   # Server private key (keep secret, never commit)
server.csr   # Server CSR (temporary)
server.crt   # Server certificate (signed by Root CA)
server.cnf   # Server config
```

---

## 1️⃣ Create Root CA

### Generate Root CA private key:
```bash
openssl genrsa -out rootCA.key 4096
```

### Create `root.cnf`:
```ini
[ req ]
distinguished_name = req_distinguished_name
x509_extensions = v3_ca
prompt = no

[ req_distinguished_name ]
C  = UK
ST = Dev
L  = Local
O  = MyOrg
OU = RootCA
CN = MyRootCA

[ v3_ca ]
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always,issuer
basicConstraints = critical, CA:true
keyUsage = critical, digitalSignature, cRLSign, keyCertSign
```

### Generate Root CA certificate:
```bash
openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 3650 \
  -out rootCA.pem -config root.cnf
```

---

## 2️⃣ Create Server Key & CSR

### Generate server private key:
```bash
openssl genrsa -out server.key 2048
```

### Create `server.cnf`:
```ini
[ req ]
default_bits       = 2048
prompt             = no
default_md         = sha256
req_extensions     = v3_req
distinguished_name = req_distinguished_name

[ req_distinguished_name ]
C  = UK
ST = Dev
L  = Local
O  = MyOrg
OU = DevServer
CN = localhost

[ v3_req ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = localhost
```

### Create Certificate Signing Request (CSR):
```bash
openssl req -new -key server.key -out server.csr -config server.cnf
```

---

## 3️⃣ Sign Server Certificate with Root CA

```bash
openssl x509 -req -in server.csr -CA rootCA.pem -CAkey rootCA.key \
  -CAcreateserial -out server.crt -days 825 -sha256 \
  -extfile server.cnf -extensions v3_req
```

---

## 4️⃣ Trust Your Root CA

### Windows
1. Double-click `rootCA.pem`  
2. Click **Install Certificate**  
3. Select **Local Machine**  
4. Place in: **Trusted Root Certification Authorities**  
5. Finish, then restart browser

### macOS
- Open **Keychain Access**  
- Import `rootCA.pem` into *System* or *System Roots*  
- Set to **Always Trust**

### Linux
```bash
sudo cp rootCA.pem /usr/local/share/ca-certificates/rootCA.crt
sudo update-ca-certificates
```

### Firefox
- Go to `about:preferences#privacy`  
- Scroll to **Certificates → View Certificates → Authorities**  
- Import `rootCA.pem`  
- Check **Trust this CA to identify websites**

---

## 5️⃣ Run a Go HTTPS Server

Example `main.go`:
```go
package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, HTTPS with my own CA!")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server running at https://localhost:8443")
	err := http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil)
	if err != nil {
		panic(err)
	}
}
```

Run:
```bash
go run main.go
```

Visit 👉 [https://localhost:8443](https://localhost:8443)  
✅ Should show **padlock, no warnings**.

---

## ⚠️ Security Notes
- **Never commit `.key` files to GitHub**  
- Add a `.gitignore` entry:
  ```
  *.key
  *.srl
  *.csr
  ```
- Root CA private key (`rootCA.key`) must be kept secure/offline in real systems  
- This setup is **for dev/learning only**, not production use  

---

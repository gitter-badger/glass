# FTP

FTP will be used as a simple preliminary protocol for payload exchange.

Before connecting: GLASS

# Sending a payload
* connect to the destination glasshouse
* PUT [UID]


The following require authentication:


# Receiving payloads
* GET [payload ID]

# Polling payloads
* ls

# Retrieve AES key
* KEY

# Retrieve AES token
* TOKEN

# Continue authentication
* USERNAME: [TOKEN]
* PASSWORD: [SALT]

import ssl

class AuthToken:
    def __init__(self, myself, appToken, appPassword, appSecret):
        self.myself = myself
        self.appClientInstanceId = appToken
        self.appPassword = appPassword
        self.appSecret = appSecret

class Myself:
    def __init__(self, pubkey):
        self.pubkey = pubkey
        self.id = "get id from pubkey"

class GlassServerConn:
    TLSPROTO = b"GLASS.TLS.0.0..."

    def __init__(self, hostname, port=443):
        self.hostname = hostname
        self.port = port
        self.auth = None
        self.__connect()

    def __connect(self):
        context = ssl.create_default_context()
        self.conn = conn = context.wrap_socket(
            socket.socket(socket.AF_INET),
            server_hostname=self.hostname
        )
        conn.connect((self.hostname, port))
        conn.sendall(self.TLSPROTO)
        self.__getAnswer("GlassHouse does not support our protocol version")

    def __getAnswer(self, errormsg):
        rep = self.conn.recv(16)
        if not rep.startswith(b"OK"):
            raise Exception(errormsg + ": `" + rep + "`.")

    def getPublicKey(self, userId):
        "GETPUBKEY [USER ID]"
        assert len(userId) == 16
        self.conn.sendall(b"GETPUBKEY " + userId + "\r\n")

    def login(self, token):
        "LOGIN [APP INSTANCE] [APP PASSWORD]"
        conn.sendall("LOGIN " + token.appClientInstanceId + b" " + token.appPassword + b"\r\n")
        self.__getAnswer("GlassHouse did not accept login")

    def instant(self, recipientAppInstanceId, payload):
        "INSTANT [RECIPIENT] [PAYLOAD SIZE]"
        conn = self.conn
        size = bytes(str(len(payload)), 'ascii')
        conn.sendall(b"INSTANT " + recipientAppInstanceId + b" " + size + b"\r\n")
        conn.sendall(payload)
        self.__getAnswer("GlassHouse did not accept payload")

    def instantInternal(self, payload):
        pass

    def pollInstant(self):
        "GETNEXTINSTANT"
        self.conn.sendall(b"GETNEXTINSTANT\r\n")
        rep = self.conn.recv(16)
        if not rep.startswith(b"OK"):
            return None
        else:

import rsa
import logging

from sleekxmpp.xmlstream import ElementBase, ET, JID, register_stanza_plugin
from sleekxmpp import Iq, Message, ClientXMPP
from sleekxmpp.xmlstream.handler import Callback
from sleekxmpp.xmlstream.matcher import StanzaPath, MatchXPath


class XEncrypted(ElementBase):
    namespace = 'glass:x:xencrypted'
    name = 'x'
    plugin_attrib = 'glass'
    sub_interfaces = interfaces = set()

class Encrypted(ElementBase):
    namespace = 'glass:x:encrypted'
    name = 'x'
    plugin_attrib = 'glass'
    sub_interfaces = interfaces = set()


class RootServer(ClientXMPP):
    """The main application running on the GlassBox."""

    def __init__(self, config):
        self.queue = []
        self.key = None
        # A mapping appId -> appInstanceIDs
        self.apps = {'': ('',)}
        # Mapping appInstanceID -> instanceKey
        self.instances = {'': ''}
        # Store here public keys
        self.pubkeys = {}
        # Admin Consoles instances
        self.admin = []

        # The V4 fingerprint of the user's public key
        self.key_fingerprint = "TODO"
        jid = 'admin@Alan' #base64_encode(self.key_fingerprint) + "@" + domain

        # connect to xmpp server with empty password
        # this will force RSA authentication instead
        ClientXMPP.__init__(self, jid, "admin")

        register_stanza_plugin(Message,  Encrypted)
        register_stanza_plugin(Message, XEncrypted)

        self.registerHandler(
          Callback("glass",
            MatchXPath('{%s}message/{glass:x:xencrypted}x' % self.default_ns),
            self.handleXEncrypted))
        self.registerHandler(
          Callback("glass",
            MatchXPath('{%s}message/{glass:x:encrypted}x' % self.default_ns),
            self.handleEncrypted))

        self.add_event_handler("session_start", self.session_start)
        #self.add_event_handler("message", self.message)

    def appID_from_instanceID(self, instanceID):
        # FIXME!
        return (k for k in self.apps if self.apps[k][0] == instanceID)[0]

    def handle(self, payload, appID):
        pass

    def handleEncrypted(self, msg):
        """Handle an encrypted payload from a contact (here, myself)."""
        if not self.key:
            return self.queue.append(("encrypted", msg))
        assert self.boundjid.bare == msg['from'].bare == msg['to'].bare
        # assert msg['type'] empty
        instanceID = msg['to'].resource
        instanceKey = self.instances[instanceID]
        appID = self.appID_from_instanceID(instanceID)
        x = msg['x']
        iv = x.attrib['iv'].bare
        data = x.bare
        payload = data # FIXME decode with instanceKey and IV
        return self.handle(payload, appID)

    def handleXEncrypted(self, msg):
        """Handle an encrypted payload from a stranger."""
        if not self.key:
            return self.queue.append(("xencrypted", msg))
        sender = msg['from'].bare
        x = msg['x']
        print("msg data", x.bare)
        # Decrypt message with user private key
        msg = rsa.decrypt(x.bare, self.key)
        # Find the appServerInstanceId
        appID, msg = msg.split(" ", 1)
        instanceID = self.apps[appID][0]
        instanceKey = self.instances[instanceID]

        # Get sender and its signature of the message
        pubkey = self.pubkeys.get(sender, None)
        if not pubkey:
            # TODO: retrieve pubkey
            return False
        try:
            rsa.verify(msg, x.attrib['sig'], pubkey)
        except rsa.pkcs1.VerificationError:
            return False

        # Empty appID means that this is the recipient app
        if not instanceID:
            # parse xml in msg
            xml = None
            return self.handle(xml, appID)

        # Forward the message to the right app
        sfrom = mto = self.boundjid.bare + "/" + instanceID
        message = self.Message(sto=mto, sfrom=sfrom)
        sm = Encrypted()
        message.xml.append(sm)
        sm.attrib['from'] = sender
        sm.attrib['orig_app'] = appID
        iv = "Random.FIXME"
        msg = "msg # TODO AES encrypt with IV above"
        sm.attrib['iv'] = iv
        sm.bare = msg
        message.send()

    def session_start(self, event):
        # send notification to all admin apps
        # to provide private key password
        pass

    def onSecretReceived(self, password):
        # use self.key for authentication w/ asymmetric key
        self.key = None
        pass #TODO

    def firstTime(self):
        # Generate Admin Console authorization
        # - Request new InstanceID, InstancePassword
        # - Generate new InstanceKey
        # - Store them in clear
        # - Add all of them amnually in the app
        pass

if __name__ == '__main__':
    # Configuration
    config = {
      'private_key_path' : "/path/to/key.pem",
      'public_key_path' : "/path/to/key.pem",
      'domain' : "localhost",
      'username' : "",
      'password' : ""
    }
    logging.basicConfig(level = logging.DEBUG)

    app = RootServer(config)
    #app.connect()
    #app.process(block=True)

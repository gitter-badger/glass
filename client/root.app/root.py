import rsa
import logging


from sleekxmpp.xmlstream import ElementBase, ET, JID, register_stanza_plugin
from sleekxmpp import Iq, Message, ClientXMPP
from sleekxmpp.xmlstream.handler import Callback
from sleekxmpp.xmlstream.matcher import StanzaPath, MatchXPath


class StrangeMessage(ElementBase):
    namespace = 'glass:main:stranger'
    name = 'c'
    plugin_attrib = 'glass'
    sub_interfaces = interfaces = set(('data', 'sig'))

#    def addField(self, name):
#        item = ET.Element('{%s}%s' % (self.namespace, name))
#        self.xml.append(item)
#        return item

class SecureMessage(ElementBase):
    namespace = 'glass:main:contact'
    name = 'c'
    plugin_attrib = 'glass'
    sub_interfaces = interfaces = set(('data', 'iv'))


class RootServer(ClientXMPP):
    """The application running on the GlassBox."""

    def __init__(self, private_key, domain):
        # A mapping appId -> appServerInstanceId
        self.apps = {"": ""}
        # Mapping appInstanceID -> instanceKey
        self.instances = {"": ""}
        # Store here public keys
        self.pubkeys = {}

        self.key = private_key
        self.domain = domain
        # The V4 fingerprint of the user's private key
        self.key_fingerprint = "TODO"
        jid = 'admin@Alan' #base64_encode(self.key_fingerprint) + "@" + domain

        # connect to xmpp server with empty password
        # this will force RSA authentication instead
        ClientXMPP.__init__(self, jid, "admin")

        self.registerHandler(
          Callback("glass",
            MatchXPath('{%s}message/{glass:main:stranger}c' % self.default_ns),
            self.handleStranger))
        register_stanza_plugin(Message, StrangeMessage)

        #self.add_event_handler("session_start", self.session_start)
        #self.add_event_handler("message", self.message)

    def handleStranger(self, msg):
        """Handle a message from a stranger."""
        sender = msg['from'].bare
        c = msg['c']
        print("msg data", c['data'].bare)
        # Decrypt message with user private key
        msg = rsa.decrypt(c['data'].bare, self.key)
        # Find the appServerInstanceId
        appID, msg = msg.split(" ", 1)

        # Get sender and its signature of the message
        pubkey = self.pubkeys.get(sender, None)
        if not pubkey:
            # TODO: retrieve pubkey
            return False
        try:
            rsa.verify(msg, c['signature'], pubkey)
        except rsa.pkcs1.VerificationError:
            return False

        # Empty appID means that this is the recipient app
        instanceID = self.apps.get(appID, "")
        instanceKey = self.instances(instanceID, "")
        if not instanceID:
            # parse xml in msg
            # handle it internally
            return

        # Forward the message to the right app
        sfrom = mto = self.boundjid.bare + "/" + instanceID
        message = self.Message(sto=mto, sfrom=sfrom)
        sm = SecureMessage()
        message.xml.append(sm)
        sm.attrib['from'] = sender
        sm.attrib['from_app'] = appID
        iv = "Random.FIXME"
        msg = "msg # TODO AES encrypt with IV above"
        sm[ 'iv' ] = iv
        sm['data'] = msg
        message.send()

    def session_start(self, event):
        # use self.key for authentication w/ asymmetric key
        pass #TODO

if __name__ == '__main__':
    # >>> with open('private.pem') as privatefile:
    #...     keydata = privatefile.read()
    #>>> pubkey = rsa.PrivateKey.load_pkcs1(keydata)

    private_key = None
    domain = "localhost"

    logging.basicConfig(level = logging.DEBUG)

    app = RootServer(private_key, domain)
    app.connect()
    app.process(block=True)

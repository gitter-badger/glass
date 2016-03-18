class Channel:
    app = ""
    id = ""
    secret = ""

class RootServer(GlassApp):
    appId = "root"
    key = None
    channels = {}

    def onSuccessfulConnection(self):
        # use self.key for authentication w/ asymmetric key
        pass #TODO

    def onReceiveMessage(self, stanza, sender):
        # decrypt message with user private key
        msg2 = self.key.decrypt(msg)
        # get sender and its signature of the message
        signature = msg.signature
        # download sender's key and check the signature (...)
        appId = msg2.appId
        if appId == this.appId:
            # the message's destination is this app
            msg3 = msg2.content
            # set correct sender
            msg3["from"] = sender
            # handle root service message
            return self.onServiceMessage(msg3)
        if appId in self.channels:
            # the destination app is registered
            # get its channel
            channel = self.channels[appId]
            # generate new initialization vector
            iv = random
            # encrypt message with channel's secret and iv
            msg3 = EAS.encrypt(msg2, channel.secret, iv)
            # send the message to the channel
            return self.send({to: channel.id, iv: iv, message: msg3})

    def onServiceMessage(self, msg):
        pass

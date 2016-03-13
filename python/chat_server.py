from glass import App, GlassApp, Payload

class ChatAppServer(GlassApp):
    def onReceiveIncomingNotification(self, payload, internal=None):
        if internal is not False:
            self.sendNotification(self, payload)
        else:
            self.sendNotification(self, payload, self.instanceId)

    def init(self):
        self.instanceId = "0123456789ABCDEF"


app = App("com.glass.example.chat", ChatAppServer)

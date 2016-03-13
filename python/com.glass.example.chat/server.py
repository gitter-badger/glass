from glass import ServeApp, GlassApp, Payload

class ChatAppServer(GlassApp):
    def onReceiveIncomingNotification(self, payload, internal=None):
        if internal is not False:
            self.sendNotification(self, payload)
        else:
            self.sendNotification(self, payload, self.clientInstance)

    def init(self):
        self.clientInstance = self.glass.newInstance()

app = ServeApp(ChatAppServer)

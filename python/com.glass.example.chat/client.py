from glass import GlassApp, Payload

class ChatApp(GlassApp):
    def onReceiveIncomingNotification(self, payload):
        print("[%s]>> %s" % (payload["__from"], payload["message"]))

    def init(self):
        while True:
            msg = input("<< ")
            if not msg:
                self.poll()
            else:
                payload = Payload(
                    __to = self.userAddress,
                    message = msg
                )
                self.sendNotification(payload, internal = "")


if __name__ == "__main__":
    instance = AppInstance(
        "user@localhost",
        "com.glass.example.chat",
        "instanceId", "instanceKey"
    )
    ChatApp(instance)

import sys

class GlassApp:
    def __init__(self, appInstance, glassClient = None):
        self.__queue = []
        self.glass = glassClient

        self.__userPubKey = None
        self.__userGlassServer = None
        self.__userId = None

        self.userAddress = appInstance.userAddress
        self.instance = appInstance

        self.init()

    def init(self):
        "Put here the initialization code."
        pass

    def onReceiveIncomingNotification(self, payload, internal=False):
        "Put here the code to handle incoming notifications."
        pass

    def poll(self):
        pass

    def sendNotification(self, payload, internal=None):
        if internal is None:
            encrypt = str(payload)
            recipient = payload['__to']
            # sign payload with self.__userPublicKey
            # encrypt with recipient's key
            self.__queue.append(("notification", recipient, encrypted))
        else:
            payload['__internal_from'] = self.__appClientInstanceId
            payload['__internal_to'] = internal
            encrypted = str(payload)
            iv = ""
            # TODO: encrypt payload with self.__instanceKey
            self.__queue.append(("internal", internal, iv, encrypted))



class Glass:
    def __init__(self, appId):
        self.appId = appId

    def newInstance(self):
        instanceId = ""
        instanceKey = ""
        return AppInstance(self.appId, instanceId, instanceKey)

class AppInstance:
    def __init__(self, userAddress, parentAppId, instanceId, instanceKey):
        self.userAddress = userAddress
        self.parentAppId = parentAppId
        self.instanceId  = instanceId
        self.instanceKey = instanceKey

    def isServer(self):
        return len(self.instanceKey) == 0

    def getServerAppInstance(self):
        pass

class Payload:
    def __init__(self, **kwargs):
        self.dict = kwargs
    def __getitem__(self, item):
        return self.dict[item]
    def __str__(self):
        return str(self.dict)

class ServeApp:
    def __init__(self, appClass):
        userId = sys.argv[1]
        appId = sys.argv[2]
        glass = Glass(userId, appId)
        appInstance = glass.getServerAppInstance()
        i = appClass(appInstance, glassClient=glass)
        # poll i indefinitely

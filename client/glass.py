from sleekxmpp.plugins.base import base_plugin

class GlassApp(base_plugin):
    def plugin_init(self):
        #self.description = "In-Band Registration"
        # self.xep = "0077"
        self.xmpp.registerHandler(
          Callback('In-Band Registration',
            MatchXPath('{%s}iq/{jabber:iq:register}query' % self.xmpp.default_ns),
            self.__handleRegistration))
        register_stanza_plugin(Iq, Registration)

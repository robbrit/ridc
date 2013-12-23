import json
from twisted.internet import reactor
from twisted.internet.defer import Deferred
from twisted.internet.protocol import Protocol, ClientCreator


class RiddickProtocol(Protocol):
    def connectionMade(self):
        self.lastDeferred = None

    def dataReceived(self, data):
        if self.lastDeferred:
            d = self.lastDeferred
            self.lastDeferred = None
            d.callback(json.loads(data))

    def sendMessage(self, message, onFinish):
        self.lastDeferred = onFinish
        self.transport.write(str("%s\n" % message))


class Riddick(object):
    def __init__(self, host, port):
        self.client = ClientCreator(reactor, RiddickProtocol)
        self.connected = False
        self.protocol = None
        self.message_queue = []
        self.sending = False
        self.client.connectTCP(host, port).addCallback(self._on_connect)

    def _on_connect(self, protocol):
        self.protocol = protocol
        self.connected = True
        self._trigger_next()

    def _send_message(self, message):
        d = Deferred()
        d.addCallback(self._trigger_next)

        if not self.protocol or self.sending:
            self.message_queue.append((message, d))
            return d

        self._trigger_send(message, d)
        return d

    def _trigger_send(self, message, deferred):
        ''' Actually send a message. '''
        self.sending = True
        self.protocol.sendMessage(message, deferred)

    def _trigger_next(self, data={}):
        self.sending = False

        if len(self.message_queue) == 0:
            return data

        msg, d = self.message_queue.pop(0)
        self._trigger_send(msg, d)
        return data

    def find_by_id(self, id):
        ''' Look up something by ID '''
        return self._send_message("GET /%s" % id)

    def find(self, field, value):
        ''' Look up something by an index '''
        return self._send_message("GET /%s/%s" % (field, value))

    def create(self, document):
        ''' Create a new document. '''
        return self._send_message("POST %s" % json.dumps(document))

    def delete(self, id):
        ''' Delete something by ID '''
        return self._send_message("DELETE /%s" % id)

    def delete_all(self, field, value):
        ''' Delete something by index '''
        return self._send_message("DELETE /%s/%s" % (field, value))

    def create_index(self, field):
        ''' Create a new index. '''
        return self._send_message("POST /indexes/%s" % field)

    def delete_index(self, field):
        ''' Delete an index. '''
        return self._send_message("DELETE /indexes/%s" % field)

    def indexes(self):
        ''' Get all the indexes. '''
        return self._send_message("GET /indexes")

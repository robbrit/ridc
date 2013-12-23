import riddick
from twisted.internet import reactor, defer


@defer.inlineCallbacks
def do_all_the_things3():
    db = riddick.Riddick("localhost", 3001)

    yield db.create_index("hello")
    yield db.create_index("sweet")

    res = yield db.indexes()

    print res

    reactor.stop()


@defer.inlineCallbacks
def do_all_the_things():
    # connect first
    db = riddick.Riddick("localhost", 3001)

    last = None
    times = 100000

    # first stick 100 elements in
    for i in range(times):
        data = {"hello": "blah", "field": str(i)}
        res = yield db.create(data)
        #print res

    # now create some indexes
    yield db.create_index("hello")

    # now add 100 more
    for i in range(times):
        data = {"hello": "blah", "field": str(i)}
        last = yield db.create(data)

    # create another index
    yield db.create_index("field")

    # show all our indexes
    res = yield db.indexes()
    print res

    # fetch something
    res = yield db.find_by_id(last["id"])
    print res

    print "Getting by index:"
    res = yield db.find("field", "8")

    if len(res) != 2:
        print "somesing is wrong"

    print res

    res = yield db.find("field", "3")
    print res

    # now delete some stuff
    res = yield db.delete_all("field", "3")
    print res
    res = yield db.delete(last["id"])
    print res

    # now try to fetch that stuff, shouldn't get it
    res = yield db.find_by_id(last["id"])
    print res

    res = yield db.find("field", "3")
    print res

    reactor.stop()

do_all_the_things()
reactor.run()

#!/usr/bin/env python3

from http.client import HTTPConnection
import json


def main():
    c = HTTPConnection("localhost", 8888)
    params = [
        dict(x=1, y=2),
        [1, 2]
    ]
    i = 0
    for p in params:
        c.request("POST", "/api/v1/rpc", json.dumps(dict(jsonrpc="2.0",
                                                         method="add",
                                                         id=i,
                                                         params=p
                                                         )))
        r = c.getresponse()
        print(json.loads(r.read()))
        i += 1


if __name__ == '__main__':
    main()

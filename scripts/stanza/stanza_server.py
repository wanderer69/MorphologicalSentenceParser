# -*- coding: utf-8 -*-
import stanza

import json
import base64
import sys
import os

from asyncio_simple_http_server import HttpServer
import asyncio

from asyncio_simple_http_server import uri_mapping

import pathlib

#if getattr(sys, "frozen", False) and hasattr(sys, "_MEIPASS"):
#    os.environ["PYMORPHY2_DICT_PATH"] = str(
#        pathlib.Path(sys._MEIPASS).joinpath("pymorphy2_dicts_ru/data")
#    )

resources_path = '/home/wanderer/stanza_resources'
try:
   resources_path = os.environ["STANZA_RESOURCES_PATH"]
except:
   pass

#stanza.download('ru') # download English model
#nlp = stanza.Pipeline('ru')
nlp = stanza.Pipeline(lang='ru', dir=resources_path, download_method=None)

def start(event):
    print("Start module")
    print(event)
    print("done")

def parse_sentence(sentence):
    # print("Parse sentence")
    # text = 'весёлый робот быстро едет в большой город на красном поезде'
    doc = nlp(sentence) # run annotation over a sentence
    s = "{}".format(doc)
    print(s)
    # message_bytes = s.encode("ascii")
    message_bytes = s.encode("utf8")
    base64_bytes = base64.b64encode(message_bytes)
    # base64_bytes = base64.b64encode(s)
    # print(s)
    return str(base64_bytes)

class WorkHandler:
    @uri_mapping('/test-get')
    def test_get(self):
        return {'a': 10}

    @uri_mapping('/phrase-translate', method='POST')
    def test_post(self, body):
      print(body)
      phrase = body.get('phrase')
      print('phrase:', phrase)
      result = parse_sentence(phrase)
      return {'result': result}

async def main():
    http_server = HttpServer()
    http_server.add_handler(WorkHandler())
    http_server.add_default_response_headers({
        'Access-Control-Allow-Origin': '*'
    })

    await http_server.start('0.0.0.0', 8888)
    await http_server.serve_forever()

s = parse_sentence("весёлый робот быстро едет в большой город на красном поезде")
# print(s)

if __name__ == '__main__':
    asyncio.run(main())

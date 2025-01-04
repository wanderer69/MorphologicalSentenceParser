# -*- coding: utf-8 -*-
from natasha import (
    Segmenter,
    MorphVocab,
    NewsEmbedding,
    NewsMorphTagger,
    NewsSyntaxParser,
    NewsNERTagger,
#    PER,
    NamesExtractor,
    Doc,
)

import json
import base64
import sys
import os

from asyncio_simple_http_server import HttpServer
import asyncio

from asyncio_simple_http_server import uri_mapping

import pathlib

if getattr(sys, "frozen", False) and hasattr(sys, "_MEIPASS"):
    os.environ["PYMORPHY2_DICT_PATH"] = str(
        pathlib.Path(sys._MEIPASS).joinpath("pymorphy2_dicts_ru/data")
    )

segmenter = Segmenter()
morph_vocab = MorphVocab()

emb = NewsEmbedding()
morph_tagger = NewsMorphTagger(emb)
syntax_parser = NewsSyntaxParser(emb)
ner_tagger = NewsNERTagger(emb)

names_extractor = NamesExtractor(morph_vocab)

# text = 'Посол Израиля на Украине Йоэль Лион признался, что пришел в шок, узнав о решении властей Львовской области объявить 2019 год годом лидера запрещенной в России Организации украинских националистов (ОУН) Степана Бандеры. Свое заявление он разместил в Twitter. «Я не могу понять, как прославление тех, кто непосредственно принимал участие в ужасных антисемитских преступлениях, помогает бороться с антисемитизмом и ксенофобией. Украина не должна забывать о преступлениях, совершенных против украинских евреев, и никоим образом не отмечать их через почитание их исполнителей», — написал дипломат. 11 декабря Львовский областной совет принял решение провозгласить 2019 год в регионе годом Степана Бандеры в связи с празднованием 110-летия со дня рождения лидера ОУН (Бандера родился 1 января 1909 года). В июле аналогичное решение принял Житомирский областной совет. В начале месяца с предложением к президенту страны Петру Порошенко вернуть Бандере звание Героя Украины обратились депутаты Верховной Рады. Парламентарии уверены, что признание Бандеры национальным героем поможет в борьбе с подрывной деятельностью против Украины в информационном поле, а также остановит «распространение мифов, созданных российской пропагандой». Степан Бандера (1909-1959) был одним из лидеров Организации украинских националистов, выступающей за создание независимого государства на территориях с украиноязычным населением. В 2010 году в период президентства Виктора Ющенко Бандера был посмертно признан Героем Украины, однако впоследствии это решение было отменено судом. '
# text = 'толстая серая кошка ест большую мышь.'
# text = 'кошка это животное из рода кошачьих'

def start(event):
    print("Start module")
    print(event)
    #    print(IntConnector.system("dir"))
    print("done")

def parse_sentence(sentence):
    # print("Parse sentence")
    # text = 'весёлый робот быстро едет в большой город на красном поезде'
    text = sentence
    doc = Doc(text)
    # print(doc)
    doc.segment(segmenter)
    # print(doc.tokens[:5])
    # print(doc.sents[:5])
    doc.tag_morph(morph_tagger)
    # print(doc.tokens[:5])
    # doc.sents[0].morph.print()

    for token in doc.tokens:
        token.lemmatize(morph_vocab)

    doc.parse_syntax(syntax_parser)
    # print(doc.tokens)
    # print(">>>>")
    # print(doc.sents[0].syntax)
    # doc.sents[0].syntax.print()

    for el in doc.sents[0].syntax.tokens:
        print(el)
    res = []
    for el in doc.tokens:
        # print(el, el.text, el.id, el.rel, el.pos)
        # print(el)
        p = {}
        p["start"] = "%s" % (el.start)
        p["stop"] = "%s" % (el.stop)
        p["text"] = el.text
        p["lemma"] = el.lemma
        p["id"] = el.id
        p["head_id"] = el.head_id
        p["rel"] = el.rel
        p["pos"] = el.pos
        pp = {}
        for k in el.feats:
            pp[k] = el.feats[k]
        ss = json.dumps(pp)
        p["feats"] = ss  # "%s" % (el.feats)
        res.append(p)
    # data = {'morph': doc.sents[0].morph, 'syntax': doc.sents[0].syntax}
    s = json.dumps(res)
    message_bytes = s.encode("ascii")
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

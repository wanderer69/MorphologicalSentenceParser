import stanza
#stanza.download('ru')
#nlp = stanza.Pipeline('ru')
#nlp = stanza.Pipeline(lang='ru', dir=args.models_dir, use_gpu=(not args.cpu))
nlp = stanza.Pipeline(lang='ru', dir='/home/wanderer/stanza_resources', download_method=None)

sentence = "Робот был построен в цеху."
sentence = "аббат - это настоятель мужского католического монастыря."
sentence = "содержание - это единство всех основных элементов целого, его свойств и связей, существующее и выражаемое в форме  и неотделимое от нее."
sentence = "исполнитель помещает список фреймов в переменную список_фреймов"
sentence = "исполнитель помещает список фреймов в переменную величину список_фреймов"
sentence = "аббат является настоятелем мужского католического монастыря."
doc = nlp(sentence) # run annotation over a sentence
#print(doc)
s = "{}".format(doc)
print(type(s))
print(s)
#
#for el in doc.sentences:
#    for el1 in el:
#        print(el1)

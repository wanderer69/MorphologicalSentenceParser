# -*- mode: python ; coding: utf-8 -*-

# pymorph data
import pymorphy2_dicts_ru
pymorph_data = pymorphy2_dicts_ru.get_path()

# Bootloader has to know the name of Python library. Pass python libname to CArchive. 
from PyInstaller.depend import bindepend
ll=bindepend.get_python_library_path()
#print("get_python_library_path", ll)
#pylib_name = os.path.basename(bindepend.get_python_library_path()) 
#print(pylib_name)

a = Analysis(
    ['natasha1.py'],
    pathex=['/home/wanderer/.local/lib/python3.12/site-packages'],
    binaries=[(ll, '.')],
    datas=[(pymorph_data, 'pymorphy2_dicts_ru/data'), 
        ('/app/emb/navec_news_v1_1B_250K_300d_100q.tar', 'natasha/data/emb'),
        ('/app/model/slovnet_morph_news_v1.tar','natasha/data/model'), 
        ('/app/model/slovnet_syntax_news_v1.tar', 'natasha/data/model'), 
        ('/app/model/slovnet_ner_news_v1.tar', 'natasha/data/model'), 
        ('/app/dict/last.txt', 'natasha/data/dict'), 
        ('/app/dict/maybe_first.txt', 'natasha/data/dict'), 
        ('/app/dict/first.txt', 'natasha/data/dict')],
    hiddenimports=[],
    hookspath=[],
    hooksconfig={},
    runtime_hooks=[],
    excludes=[],
    noarchive=False,
    optimize=0,
)
pyz = PYZ(a.pure)

exe = EXE(
    pyz,
    a.scripts,
    a.binaries,
    a.datas,
    [],
    name='natasha1',
    debug=False,
    bootloader_ignore_signals=False,
    strip=False,
    upx=True,
    upx_exclude=[],
    runtime_tmpdir=None,
    console=True,
    disable_windowed_traceback=False,
    argv_emulation=False,
    target_arch=None,
    codesign_identity=None,
    entitlements_file=None,
)

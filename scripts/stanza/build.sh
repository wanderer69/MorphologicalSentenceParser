pyinstaller \
  --name stanza_server \
  --onefile stanza_server.py \
  --paths "$(python3 -m site --user-site)"

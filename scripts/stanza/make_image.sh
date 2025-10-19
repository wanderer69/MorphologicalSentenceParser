mkdir ./data
cd ./data
mkdir ./stanza_resources
cd ./stanza_resources
mkdir ./ru
cd ..
wget https://github.com/stanfordnlp/stanza-resources/archive/refs/heads/main.zip
unzip main.zip
cd ./stanza-resources-main
cp ./resources_1.11.0.json ../stanza_resources/resources.json
cd ..
rm -r ./stanza-resources-main
rm main.zip
git lfs install
git clone https://huggingface.co/stanfordnlp/stanza-ru
cd ./stanza-ru
cd ./models
rm default.zip
cp -r ./backward_charlm  ../../stanza_resources/ru                           
cp -r ./coref  ../../stanza_resources/ru 
cp -r ./depparse  ../../stanza_resources/ru 
cp -r ./forward_charlm  ../../stanza_resources/ru 
cp -r ./lemma ../../stanza_resources/ru  
cp -r ./ner ../stanza_resources/ru  
cp -r ./pos ../../stanza_resources/ru  
cp -r ./pretrain ../../stanza_resources/ru  
cp -r ./tokenize ../../stanza_resources/ru
cd ../..
rm -r -f ./stanza-ru
cd ..
#docker buildx build -t stanza_server:latest --progress=plain . &> build.log
					  
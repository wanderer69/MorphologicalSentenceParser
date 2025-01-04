mkdir ./data
cd ./data
wget https://github.com/natasha/natasha/archive/refs/heads/master.zip
unzip master.zip
cd ./natasha-master/natasha
cp -r ./data ../../../
cd ../..
rm -r ./natasha-master
rm  __init__.py
rm master.zip
docker buildx build -t natasha1:latest --progress=plain . &> build.log

export RULES_FILE_NAME=../../../SemanticNet/data/rules.script
cd ../cmd/service
if [ ! -f ./service ]
then
    echo "Service not found"
    go build
fi
./service

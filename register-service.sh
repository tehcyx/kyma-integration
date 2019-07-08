#!/bin/bash
curl -X POST -d @service.json https://gateway.kyma.local:30019/github-test/v1/metadata/services --cert bin/.key/server.crt --key bin/.key/client.key -k

curl https://gateway.kyma.local:30019/github-test/v1/metadata/api.yaml --cert bin/.key/server.crt --key bin/.key/client.key -k

curl https://gateway.kyma.local:30019/github-test/v1/metadata/services/654dd210-d884-4cc3-b8b4-5a95578874f4 --cert bin/.key/server.crt --key bin/.key/client.key -k
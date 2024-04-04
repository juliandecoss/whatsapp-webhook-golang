package:
	sls package

package-local:
	sls package --stage local --package go-example-artifacts

local-run:
	sls invoke local --function whatsapp-webhook-golang --path ./test/event_body.json

install-local:
	npm install --save-dev serverless-dotenv-plugin	
	npm install --save-dev serverless-localstack
	npm i --save-dev serverless-go-plugin

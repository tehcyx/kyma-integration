# how to run locally

* use [kyma-cli](https://github.com/kyma-incubator/kyma-cli), to install kyma.
* find out nodeport of application via `kubectl -n kyma-system get svc application-connector-nginx-ingress-controller -o 'jsonpath={.spec.ports[?(@.port==443)].nodePort}'`
* change port in the url in file `kis_register.go:161` to the port, that you got from the previous command
* run `make`
* open browser and connect to your local kyma cluster to create an application [https://console.kyma.local/home/apps](https://console.kyma.local/home/apps)
* currently pick application name as 'github-test', comment 'test'
* click on connect application and copy the URL with the token
* run the application via `./bin/app`
    * output should show the application running on port 8080
* in the browser browse to http://localhost:8080/connect?url=<TOKEN_URL> and append the token url you copied in the step before as part of the url
* browser should show successful connection, application log should show that the client got updated with the TLS certificate and started listening on 8443 as well
* navigate to [http://localhost:8080/register-service](http://localhost:8080/register-service) to register the applications API with kyma
* refreshing this url [https://console.kyma.local/home/apps/details/github-test](https://console.kyma.local/home/apps/details/github-test) should show **Daniel's Service** in the *Provided Services & Events* section of the page with a green checkmark on the API


# what's currently working

* Automatically generate certificates, public and private keys.
* Connect to kyma
* Regenaration of certificate somehow broken as it will empty out the cert file
    * if this is fixed it would allow for auto-restart of TLS listener as this is already implemented
* register of API is working
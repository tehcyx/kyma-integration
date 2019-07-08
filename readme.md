# how to run locally

* use [kyma-cli](https://github.com/kyma-incubator/kyma-cli), to install kyma.
* run `make`
* open browser and connect to your local kyma cluster to create an application [https://console.kyma.local/home/cmf-apps](https://console.kyma.local/home/cmf-apps)
* add a new application
* click on connect application and copy the URL with the token
* run the application via `./bin/kis`
    * output should show the application running on port 8080:
        ```
        2019/07/08 10:25:18 🔓 Listening on 0.0.0.0:8080
        ```
* in the browser browse to [http://localhost:8080/connector/connect?url=<TOKEN_URL>](http://localhost:8080/connector/connect?url=<TOKEN_URL>) and append the token url you copied in the step before as part of the url
    * browser should show successful connection, application log should show that the client got updated with the TLS certificate and started listening on 8443 as well
        ```
        2019/07/08 10:25:59 handler started
        2019/07/08 10:26:04 updating http client with certificate
        2019/07/08 10:26:04 handler ended
        2019/07/08 10:26:04 🔐 Listening on 0.0.0.0:8443
        ```
* navigate to [http://localhost:8080/connector/register-service](http://localhost:8080/connector/register-service) to register the applications API with kyma
    * browser should show the service id if registered succesfully
* Your created application in [https://console.kyma.local/home/cmf-apps](https://console.kyma.local/home/cmf-apps) should show **Daniel's Service** in the *Provided Services & Events* section of the page with a green checkmark on the API


# what's currently working

* Automatically generate certificates, public and private keys.
* Connect to kyma
* Regenaration of certificate might be broken (? i think i fixed it in an earlier patch)
    * needs certificate renewal functionality for 92 days valid certificates
* register of API is working, events are missing but should be the same flow with a different url and schema


# Kyma and Kyma-integration on Minikube

* Install kyma, e.g. with [kyma-cli](https://github.com/kyma-incubator/kyma-cli)

> So to use an image without uploading it, you can follow these steps:
> 
> * Set the environment variables with eval $(minikube docker-env)
> * Build the image with the Docker daemon of Minikube (eg docker build -t my-image .)
> * Set the image in the pod spec like the build tag (eg my-image)
> * Set the imagePullPolicy to Never, otherwise Kubernetes will try to download the image.
# How To Run Locally

* Use [kyma-cli](https://github.com/kyma-incubator/kyma-cli), to install kyma.
* Run `make run`
    * Output should show the application running on port 8080:
        ```
        time="2020-01-21T19:51:39-08:00" level=info msg="🔓 Listening on 0.0.0.0:8080"
        ```
    * A `config.yaml` is created in the directory where the binary is, if you don't specify a different path in the `CONFIG_DIR` environment variable.
* Open browser and connect to your local kyma cluster to create an application [https://console.kyma.local/home/cmf-apps](https://console.kyma.local/home/cmf-apps)
* Add a new application
* Click on connect application and copy the URL with the token
* In the browser browse to [http://localhost:8080/](http://localhost:8080/)
* Use the copied url from the created application on Kyma, put it in the text box and hit submit
    * This will connect the application in Kyma as well as registering API and Events with the application.
    * Double check in the Kyma application details, that after a refresh the application is listed as "Serving" and has APIs and Events listed in the *Provided Services & Events* section of the page.

# What's Currently Working

* Automatically generate certificates, public and private keys.
* Connect to Kyma and register API and Events automatically.

# Kyma and Kyma-integration on Minikube

* Install kyma, e.g. with [kyma-cli](https://github.com/kyma-incubator/kyma-cli)

> To use an image without uploading it to a docker registry, you can follow these steps:
> 
> * Set the environment variables with eval $(minikube docker-env)
> * Build the image with the Docker daemon of Minikube (eg docker build -t my-image .)
> * Set the image in the pod spec like the build tag (eg my-image)
> * Set the imagePullPolicy to Never, otherwise Kubernetes will try to download the image.

# How To Run On k8s/kind

* Install Kyma on minikube or use a remote Kyma.
* Run `make setupkind` to initialize a kind cluster just for this application.
* Run `make docker` to create a docker binary
* Run `docker build -f build/package/Dockerfile -t kyma-integration:$(shell cat VERSION)` to create a docker container and tag it with the current version.
* Run `kind load docker-image kyma-integration:$(shell cat VERSION) --name "config-map"` to load the image into kind.
* Run `kubectl apply -f deploy-k8s.yaml` to deploy the service.
* Use `kubectl ...` to retrieve NodePort of application and access the service.
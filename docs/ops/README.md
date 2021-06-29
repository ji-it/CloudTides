# Ops

## Docker

The UI and server are built into Docker images in CD process. Two [Dockerfiles](https://docs.docker.com/engine/reference/builder/) are located in `ui` and `tides-server` folders. They make use of [multi-stage builds](https://docs.docker.com/develop/develop-images/multistage-build/) to reduce image size.

## CI & CD Configuration

Reference material:
- [Introductory tutorial](http://www.ruanyifeng.com/blog/2019/09/getting-started-with-github-actions.html)
- [Syntax manual](https://docs.github.com/cn/actions)

Useful third-party actions:
- [scp](https://github.com/marketplace/actions/scp-files)
- [remote ssh](https://github.com/marketplace/actions/remote-ssh-commands)
- [publish to docker hub](https://github.com/elgohr/Publish-Docker-Github-Action)
  
Future improvements:
- The test stage for UI is missing as `npm test` always gives some errors.

## Kubernetes

CloudTides is deployed on [Kubernetes](https://kubernetes.io/) of Aliyun as three microservices, i.e. UI, server and database. The Kubernetes deployment file is located in `ui` and `tides-server` folders.

To get started, refer to [Build Kubernetes Cluster inside China](https://github.com/scienterprise/CloudTides/wiki/Build-Kubernetes-Cluster-inside-China).

The CI & CD pipeline automatically deploys UI and server on Kubernetes.

To inspect these services, `ssh` to the master node, use `kubectl` command.
```
kubectl get all/pod/service/deployment
kubectl apply -f <deployment yaml file>
kubectl describe pod/service/deployment <target-name>
kubectl logs <pod-name>
kubectl delete pod/service/deployment <target-name>
kubectl delete -f <deployment yaml file>
```

Future improvements:
- Current Kubernetes yaml file and CD workflow uses personal Docker account. Better to create a group Docker account.

## PostgreSQL

For deployment of PostgreSQL, refer to [PostgreSQL service](https://severalnines.com/database-blog/using-kubernetes-deploy-postgresql). The k8s yaml files are located in `postgres` folder. Customize your database username, password, db name, service port in those files. Then apply following commands:
```
kubectl create -f postgres-configmap.yml
kubectl create -f postgres-storage.yml
kubectl create -f postgres-deployment.yml
kubectl create -f postgres-service.yml
```

## Credentials

The credentials for backend is persistently stored in k8s [secrets](https://kubernetes.io/docs/concepts/configuration/secret/). The server container loads the secret data when it is deployed on k8s:

```
envFrom:
    - secretRef:
        name: cloudtides-secret
```

Ensure the security of credentials.

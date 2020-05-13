# Two Microservices, reverse and random

For local development use kind:

1. `cd scripts && ./start_kind.sh`
2. run the helper script `./load_kind.sh yourtag`

This will build and tag docker images for both services, load them into kind, update the image tags in the deployment files and deploy the services.

The deployment files are in `/kube`

## Deployment to GCP

1. Use GKE to bring up a cluster supporting load balancers
2. Modify the 'load_kind.sh' script to push to docker hub instead of load an image into kind
3. Deploy `kube`, and retrieve the IP of the loadbalancer for the `random` sevice.

See [diagram](diagram.png)

```bash
curl  --header "Content-Type: application/json" \
  --request POST \
  --data '{"message":"abcdef123!@£"}' \
  http://<random_lb_ip>/api
```

## CI/CD

1. Use something like CircleCI to recieve webhooks when a push/merge to the repository is made (perhaps limit it to one branch)
2. When a webhook is recieved the CI/CD service should: build, tag, push (use modified version of `load_kind.sh` with the commit SHA as the tag)
3. Use a CI/CD service account with appropriate RBAC for this project to update the deployment manifests and roll out the deployment from the CI/CD platform using a script.

As the`.spec.strategy.type` on the deployments is  set to `RollingUpdate` with `maxUnavailable` and `maxSurge` at 25% (by default) and there are readiness / liveness probes enabled, kubernetes will gradually roll out the changes and stop the rollout if a  probe fails. If a rollout is stuck for a duration longer than `progressDeadlineSeconds` a condition is set with `Reason=ProgressDeadlineExceeded`, `Status=False` and `Type=Progressing`. The script can sleep, watch for  conditions. It could exit on success and if the faulure condition is observed call `kubectl rollout undo`.

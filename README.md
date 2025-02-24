# rollouts-autoscaling-example
Example for using Argo Rollouts with autoscaling 


## Example 01

kubectl argo rollouts get rollout 01-baseline-bg

kubectl argo rollouts set image 01-baseline-bg baseline-demo=docker.io/kostiscodefresh/summer-of-k8s-app:v2

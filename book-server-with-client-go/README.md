# Book-Server with k8s.io/client-go...

## Commands to get the dependecies

- `glide init`
- `glide update`
- `glide install`
  
## Example commands to run

- The following command will create a deployment nammed 'book-server-deployment',

  `go run main.go create deploy`

- The following command will create a service nammed 'book-server-service',

  `go run main.go create svc`

- The following command will delete if there exist a deployment nammed 'book-server-deployment'
and a service nammed 'book-server-service',

  `go run main.go delete`

- If we want to tell the kube config file path, we can do that using a flag nammed 'kubeconfig' in 
all of the above commands, i.e.

  `go run main.go create svc --kubeconfig=$HOME/.kube/config`
  

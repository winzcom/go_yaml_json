package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
)

var yaml2 string = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app-deployment
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: my-app-container
        image: my-app-image
        ports:
        - containerPort: 80
        env:
        - name: DATABASE_HOST
          value: db.example.com
        - name: API_KEY
          valueFrom:
            secretKeyRef:
              name: my-secret
              key: api-key`

var yaml string = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: example-deployment
  label: mikl
  now: you know
spec:
  replicas: 3
  selector:
    matchLabels:
      app: example
  template:
    metadata:
      labels:
        app: example
    spec:
      containers:
      - name: example-container
        image: example-image
        use: jk
      - name: example-container
        image: example-image
        port:
        - hello: world
          what: amaing
        - per: not
`

func main() {
	str := strings.NewReader(yaml2)

	bfr := bufio.NewReader(str)

	builder := BuildJSON(bfr)

	to_json, _ := json.Marshal(builder)

	fmt.Println("what is builder ", string(to_json))

}

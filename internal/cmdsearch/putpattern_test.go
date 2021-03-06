package cmdsearch

var putPatternCases = []test{
	{
		name: "put pattern single setter",
		args: []string{"--by-value", "3", "--put-pattern", "${replicas}"},
		input: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3
 `,
		inputKptfile: `apiVersion: v1alpha1
kind: Kptfile
openAPI:
  definitions:
    io.k8s.cli.setters.replicas:
      x-k8s-cli:
        setter:
          name: replicas
          value: "3"`,
		out: `${baseDir}/
matched 1 field(s)
${filePath}:  spec.replicas: 3
`,
		expectedResources: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3 # {"$kpt-set":"${replicas}"}
 `,
	},
	{
		name: "put pattern group of setters",
		args: []string{"--by-value", "nginx-deployment", "--put-pattern", "${image}-${kind}"},
		input: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3
 `,
		inputKptfile: `apiVersion: v1alpha1
kind: Kptfile
openAPI:
  definitions:
    io.k8s.cli.setters.image:
      x-k8s-cli:
        setter:
          name: image
          value: "nginx"
    io.k8s.cli.setters.kind:
      x-k8s-cli:
        setter:
          name: kind
          value: "deployment"`,
		out: `${baseDir}/
matched 1 field(s)
${filePath}:  metadata.name: nginx-deployment
`,
		expectedResources: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment # {"$kpt-set":"${image}-${kind}"}
spec:
  replicas: 3
 `,
	},
	{
		name: "put pattern by regex",
		args: []string{"--by-value-regex", "my-project-*", "--put-pattern", "${project}-*"},
		input: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-project-deployment
  namespace: my-project-namespace
spec:
  replicas: 3
 `,
		inputKptfile: `apiVersion: v1alpha1
kind: Kptfile
openAPI:
  definitions:
    io.k8s.cli.setters.project:
      x-k8s-cli:
        setter:
          name: project
          value: "my-project"`,
		out: `${baseDir}/
matched 2 field(s)
${filePath}:  metadata.name: my-project-deployment
${filePath}:  metadata.namespace: my-project-namespace
`,
		expectedResources: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-project-deployment # {"$kpt-set":"${project}-deployment"}
  namespace: my-project-namespace # {"$kpt-set":"${project}-namespace"}
spec:
  replicas: 3
 `,
	},
	{
		name: "put pattern by value",
		args: []string{"--by-value", "dev/my-project/nginx", "--put-pattern", "${env}/${project}/${name}"},
		input: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dev/my-project/nginx
spec:
  replicas: 3
 `,
		inputKptfile: `apiVersion: v1alpha1
kind: Kptfile
openAPI:
  definitions:
    io.k8s.cli.setters.project:
      x-k8s-cli:
        setter:
          name: project
          value: "my-project"
    io.k8s.cli.setters.env:
      x-k8s-cli:
        setter:
          name: env
          value: "dev"
    io.k8s.cli.setters.name:
      x-k8s-cli:
        setter:
          name: name
          value: "nginx"
    io.k8s.cli.setters.namespace:
      x-k8s-cli:
        setter:
          name: namespace
          value: "my-space"`,
		out: `${baseDir}/
matched 1 field(s)
${filePath}:  metadata.name: dev/my-project/nginx
`,
		expectedResources: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dev/my-project/nginx # {"$kpt-set":"${env}/${project}/${name}"}
spec:
  replicas: 3
 `,
	},
	{
		name: "put pattern error",
		args: []string{"--by-value", "nginx-deployment", "--put-pattern", "${image}-${tag}"},
		input: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3
 `,
		inputKptfile: `apiVersion: v1alpha1
kind: Kptfile
openAPI:
  definitions:
    io.k8s.cli.setters.image:
      x-k8s-cli:
        setter:
          name: image
          value: "nginx"
    io.k8s.cli.setters.kind:
      x-k8s-cli:
        setter:
          name: kind
          value: "deployment"`,
		expectedResources: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3
 `,
		errMsg: `setter "tag" doesn't exist, please create setter definition and try again`,
	},
	{
		name: "put pattern list-values error",
		args: []string{"--by-value", "3", "--put-pattern", "${replicas-list}"},
		input: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3
 `,
		inputKptfile: `apiVersion: v1alpha1
kind: Kptfile
openAPI:
  definitions:
    io.k8s.cli.setters.replicas-list:
      x-k8s-cli:
        setter:
          name: replicas-list
          value: ""
          listValues: 
           - "1"
           - "2"`,
		expectedResources: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3
 `,
		errMsg: `setter pattern should not refer to array type setters: "replicas-list"`,
	},
}

<h1 align="center" style="border-bottom: none">
    <a href="https://nimtechnology.com/2023/07/02/jougan-project/" target="_blank"><img alt="JouGan" width="120px" src="https://cdn131.picsart.com/329673960061211.png"></a><br>JouGan
</h1>

<p align="center">Monitoring the speed and volume of Kubernetes in a real-life scenario.</p>

## Introduction
JouGan can download the file you request. I calculate and report:   
- Time taken to download the file (seconds)
- Download speed (MB/s)
- Time taken to save the file (seconds)
- Save speed (MB/s)
- Time taken to delete the file (seconds)
- Delete speed (MB/s)

## Install Jougan on K8s   

Add repository:   
```shell
helm repo add jougan https://underndog.github.io/jougan
```

Install chart:
```shell
helm install my-jougan jougan/jougan --version x.x.x
```

### value file   

Example:   
#### Measure Disk Speed any file from Download URL
```yaml
envVars:
  DOWNLOAD_URL: "https://files.testfile.org/PDF/10MB-TESTFILE.ORG.pdf"
  SAVE_TO_LOCATION: /app/downloaded/dynamicSize.bin
nodeSelector:
  kubernetes.io/os: linux
podAnnotations:
  prometheus.io/scrape: 'true'
  prometheus.io/path: '/metrics'
  prometheus.io/port: '1994'
volumes:
  - name: file-service
    persistentVolumeClaim:
    claimName: pvc-file-service-smb-1
volumeMounts:
  - mountPath: /app/downloaded
    name: file-service
```

#### Measure Disk Speed any file on S3 (AWS)   
```yaml
envVars:
  AWS_REGION: us-west-2
  DOWNLOAD_FROM_S3_BUCKET: "ahihi-09262023"
  DOWNLOAD_FROM_S3_KEY: "10MB-TESTFILE.ORG.pdf"
  SAVE_TO_LOCATION: "/app/downloaded/dynamicSize.bin"
  AWS_ACCESS_KEY_ID: "XXXXXGKBQ65KXXXXXX"
  AWS_SECRET_ACCESS_KEY: "xxxxxxx//1dSxxxxxxxJ7nkIrxxxxxxx"
nodeSelector:
  kubernetes.io/os: linux
podAnnotations:
  prometheus.io/scrape: 'true'
  prometheus.io/path: '/metrics'
  prometheus.io/port: '1994'
volumes:
  - name: file-service
    persistentVolumeClaim:
    claimName: pvc-file-service-smb-1
volumeMounts:
  - mountPath: /app/downloaded
    name: file-service
```

### Yaml deployment:   
You install quickly Jougan

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io/instance: jougan
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/name: jougan
    app.kubernetes.io/version: v0.0.2
    helm.sh/chart: jougan-0.1.1
  name: jougan
  namespace: mdaas-engines-dev
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/instance: jougan
      app.kubernetes.io/name: jougan
  template:
    metadata:
      annotations:
        prometheus.io/path: /metrics
        prometheus.io/port: '1994'
        prometheus.io/scrape: 'true'
      labels:
        app.kubernetes.io/instance: jougan
        app.kubernetes.io/managed-by: Helm
        app.kubernetes.io/name: jougan
        app.kubernetes.io/version: v0.0.2
        helm.sh/chart: jougan-0.1.1
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: app.kubernetes.io/name
                    operator: In
                    values:
                      - jougan
              topologyKey: kubernetes.io/hostname
      containers:
        - env:
            - name: AWS_ACCESS_KEY_ID
              value: XXXXXGKBQ65KXXXXXX
            - name: AWS_REGION
              value: us-west-2
            - name: AWS_SECRET_ACCESS_KEY
              value: xxxxxxx//1dSxxxxxxxJ7nkIrxxxxxxx
            - name: DEBUG_LOG
              value: 'false'
            - name: DOWNLOAD_FROM_S3_BUCKET
              value: ahihi-09262023
            - name: DOWNLOAD_FROM_S3_KEY
              value: 200MB-TESTFILE.ORG.pdf
            - name: SAVE_TO_LOCATION
              value: /app/downloaded/dynamicSize.bin
          image: 'quay.io/underndog/jougan'
          imagePullPolicy: IfNotPresent
          livenessProbe:
            httpGet:
              path: /
              port: http
          name: jougan
          ports:
            - containerPort: 1994
              name: http
              protocol: TCP
          readinessProbe:
            httpGet:
              path: /
              port: http
          resources: {}
          securityContext: {}
          volumeMounts:
            - mountPath: /app/downloaded
              name: file-service
      nodeSelector:
        kubernetes.io/os: linux
      securityContext: {}
      serviceAccountName: jougan
      volumes:
        - name: file-service
          persistentVolumeClaim:
            claimName: pvc-file-service-smb-1
```

| Environment Variable | Description | Value Example | Purpose |
| --- | --- | --- | --- |
| `AWS_ACCESS_KEY_ID` | Stores the AWS Access Key ID, which is part of the credentials used to authenticate requests to AWS services. | `XXXXXGKBQ65KXXXXXX` | Identifies the IAM user or role making the request. |
| `AWS_REGION` | Specifies the AWS region where your operations will take place. | `us-west-2` | Ensures that the application interacts with AWS services in the correct region. |
| `AWS_SECRET_ACCESS_KEY` | Contains the AWS Secret Access Key, the second part of the credentials for authenticating requests. | `xxxxxxx//1dSxxxxxxxJ7nkIrxxxxxxx` | Works with the AWS Access Key ID to authenticate the user or service making the request. |
| `DEBUG_LOG` | Indicates whether debug logging is enabled or not. | `'false'` | Controls whether the application should produce detailed logs for debugging. |
| `DOWNLOAD_FROM_S3_BUCKET` | Specifies the name of the S3 bucket from which a file will be downloaded. | `ahihi-09262023` | Tells the application which S3 bucket to access for downloading the required file. |
| `DOWNLOAD_FROM_S3_KEY` | Contains the key (or file path) of the object within the S3 bucket that needs to be downloaded. | `200MB-TESTFILE.ORG.pdf` | Identifies the exact file within the S3 bucket that the application should download. |
| `SAVE_TO_LOCATION` | Indicates the local file path where the downloaded file from S3 should be saved. | `/app/downloaded/dynamicSize.bin` | Specifies the destination directory and filename where the downloaded content will be stored locally. |
| `DOWNLOAD_TYPE` | **(Optional)** Defines the method used to download the file from S3. Options include: Default (using a Re-Signed URL only for measuring download) or AWS-S3-SDK (using AWS's official S3 SDK for measuring download). | `Default` or `AWS-S3-SDK` | Determines whether to use a Re-Signed URL or AWS's S3 SDK for downloading the file. |
| `UPLOAD_FILE_TO_S3` | **(Optional)** Measure file uploads to S3. | `false` or `true`. Default is `false` | `UPLOAD_FILE_TO_S3` lets you enable or disable the measurement of file uploads to S3, including upload speed and time. |
| `PART_SIZE_MB` | **(Optional)** PART\_SIZE is an environment variable that specifies the size of each chunk (part) of a file to be downloaded from S3. The size is expressed in megabytes (MB). | `5` | PART\_SIZE allows you to control the size of each download chunk in MB. This helps optimize download speed and efficiency by enabling the concurrent download of smaller parts of a large file. If not set, the entire file is downloaded in a single request. |
| `SHA-256-CHECKSUM` | **(Optional)** Include the `SHA-256-CHECKSUM` of the file; Jougan will verify it after download. | `dfb81a5c3f3ae4cd6bc469390e3668f2a8f3e8546f1864719673da0d8b058237` | Ensure your file remains unchanged after downloading. |
| `RANDOM_FILENAME_TO_SAVE_LOCAL` | **(Optional)** When true, Jougan will add a extra random string to the file name before saving it locally. | `false` or `true`. Default is `false` | Creating multiple pods or jougans to download the same file from S3 can lead to issues with saving and deleting the file locally. |



## Grafana

Links: https://grafana.com/grafana/dashboards/20013-jougan-measure-disk-speed/

<a href="https://nimtechnology.com/2023/07/02/jougan-project/" target="_blank"><img alt="JouGan" src="https://grafana.com/api/dashboards/20013/images/15212/image"></a>

## Create Helm Chart
helm package ./helm-chart/jougan --destination ./helm-chart/   
helm repo index . --url https://underndog.github.io/jougan   

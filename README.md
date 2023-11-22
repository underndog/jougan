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
helm repo add jougan https://mrnim94.github.io/jougan
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


## Create Helm Chart
helm package ./helm-chart/jougan --destination ./helm-chart/
helm repo index . --url https://mrnim94.github.io/jougan
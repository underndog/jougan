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


## Create Helm Chart
helm package ./helm-chart/jougan --destination ./helm-chart/
helm repo index . --url https://mrnim94.github.io/jougan

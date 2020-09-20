![image](https://user-images.githubusercontent.com/40491079/93712746-651c3100-fb60-11ea-89e9-4c207f7db7ef.png)

# kube-ns-cleaner

## description

Simple cron job for automatic deletion kubernetes namespaces via bash script.

## usage

Automatic deletion of kubernetes namespaces that exist more than 10d.
Every day at 00:00:00.

You can modify script or cron job and set another values.

`kubectl apply -f ~/kube-ns-cleaner`

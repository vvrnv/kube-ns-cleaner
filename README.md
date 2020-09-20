![repo logo] 
(https://github.com/vvoronov2/kube-ns-cleaner/blob/master/images/kube-ns-cleaner.png?raw=true)

# kube-ns-cleaner

# description

Simple cron job for automatic deletion kubernetes namespaces via bash script.

# usage

Automatic deletion of kubernetes namespaces that exist more than 10d.
Every day at 00:00:00.

You can modify script or cron job and set another values.

`kubectl apply -f ~/kube-ns-cleaner`
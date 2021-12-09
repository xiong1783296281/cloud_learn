#K8S搭建手册-Kubeadmin版
##安装要求
1.一台及以上机器,linux系统(ubuntu-server 20.04LTS)  
2.硬件要求: 2CPU 2GB内存   
3.机器间内网可访问且通外网  
4.关闭防火墙，swap

###关闭swap

```text
临时关闭
swapoff -a

永久关闭
vim /etc/fstab
注释掉最后一行(swap)
```

##准备工作
所有机器进行update
```shell
sudo apt-get update
```
进行时间同步
```shell
sudo apt-get install ntpdate
```
将时间同步加入cron任务
```shell
echo '*/5 * * * * ntpdate cn.pool.ntp.org' >>/var/spool/cron/root
```
安装相关依赖
```shell
sudo apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common
```
添加docker的官方GPG密钥
```shell
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
```
添加docker的阿里云仓库
```shell
sudo add-apt-repository \   
   "deb [arch=amd64] https://mirrors.ustc.edu.cn/docker-ce/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
```
再一次更新系统以及安装docker
```shell
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io
```

##安装k8s
install kubeadm,kubelet,kubectl
```shell
cat <<EOF > /etc/apt/sources.list.d/kubernetes.list
 deb https://mirrors.aliyun.com/kubernetes/apt kubernetes-xenial main
 EOF
```
```shell
curl -s https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | sudo apt-key add
```
```shell
sudo apt-get update
sudo apt-get install -y kubelet=1.19.15-00 kubeadm=1.19.15-00 kubectl=1.19.15-00
sudo apt-mark hold kubelet kubeadm kubectl
```
将kubelet加入开机自启服务列表
```shell
systemctl enable kubelet
```
kube init(仅master节点)
```shell
kubeadm init \   
 --image-repository registry.aliyuncs.com/google_containers \   #国内镜像源
 --kubernetes-version v1.19.15 \
 --apiserver-advertise-address=192.168.34.2 \     #master节点内网ip
 --service-cidr=10.96.0.0/12 \                    #svc网络段
 --pod-network-cidr=10.244.0.0/16                 #pod网络段
```
注: svc和pod网段不能相同
根据init结束后的提示  执行下面命令
```shell
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config
```
添加字节点进入集群(在子机器执行)
```shell
kubeadm join 192.168.34.2:6443 --token 1sjpqf.95tgknvlo99bdj9w  \     
  --discovery-token-ca-cert-hash sha256:0301a4154e8da88419eec721f2e0401118dad3f61f08064ea8b02903c8be350e
```

##安装calico
下载calico配置文件
```shell
curl https://docs.projectcalico.org/manifests/calico.yaml -O
```
修改calico配置文件pod网段
```shell
vim calico.yaml
```
配置文件下面这个配置默认是被注释的
```shell
# - name: CALICO_IPV4POOL_CIDR
#   value: "192.168.0.0/16"
```
将value值改为kubeadm初始化时，指定的pod网络段
```shell
- name: CALICO_IPV4POOL_CIDR
  value: "10.244.0.0/16"
```

检验是否安装完成
```shell
root@k8s-master:~# kubectl get pods -n kube-system -o wide
NAME                                       READY   STATUS    RESTARTS   AGE   IP               NODE         NOMINATED NODE   READINESS GATES
calico-kube-controllers-659bd7879c-2wsgd   1/1     Running   2          42d   10.244.36.67     k8s-node1    <none>           <none>
calico-node-7m5xr                          1/1     Running   1          42d   192.168.34.4     k8s-node2    <none>           <none>
calico-node-f7957                          1/1     Running   2          42d   192.168.34.3     k8s-node1    <none>           <none>
calico-node-nhn4d                          1/1     Running   3          42d   192.168.34.2     k8s-master   <none>           <none>
coredns-6d56c8448f-bd942                   1/1     Running   3          42d   10.244.235.198   k8s-master   <none>           <none>
coredns-6d56c8448f-zxz92                   1/1     Running   3          42d   10.244.235.197   k8s-master   <none>           <none>
etcd-k8s-master                            1/1     Running   3          42d   192.168.34.2     k8s-master   <none>           <none>
kube-apiserver-k8s-master                  1/1     Running   3          42d   192.168.34.2     k8s-master   <none>           <none>
kube-controller-manager-k8s-master         1/1     Running   3          42d   192.168.34.2     k8s-master   <none>           <none>
kube-proxy-2vlj7                           1/1     Running   2          42d   192.168.34.4     k8s-node2    <none>           <none>
kube-proxy-4z76k                           1/1     Running   2          42d   192.168.34.3     k8s-node1    <none>           <none>
kube-proxy-dkd27                           1/1     Running   3          42d   192.168.34.2     k8s-master   <none>           <none>
kube-scheduler-k8s-master                  1/1     Running   3          42d   192.168.34.2     k8s-master   <none>           <none>

```
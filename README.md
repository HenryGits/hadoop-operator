# hadoop-operator
Hadoop Operator on kubernetes


3. netrc 文件

   根据官网文档，设置git私有仓库用户名密码，防止 Docker 镜像编译的时候需要手动输入用户名密码

   [Frequently Asked Questions (FAQ) - The Go Programming Language](https://golang.org/doc/faq#git_https)


4. 安装 Operator-SDK

[Installation | Operator SDK](https://sdk.operatorframework.io/docs/installation/)

```shell
export ARCH=$(case $(uname -m) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n $(uname -m) ;; esac)
export OS=$(uname | awk '{print tolower($0)}')
export OPERATOR_SDK_DL_URL=https://github.com/operator-framework/operator-sdk/releases/download/v1.11.0
curl -LO ${OPERATOR_SDK_DL_URL}/operator-sdk_${OS}_${ARCH}
chmod +x operator-sdk_${OS}_${ARCH} && mv operator-sdk_${OS}_${ARCH} /usr/local/bin/operator-sdk
echo 'source <(operator-sdk completion bash)' >>~/.bashrc
operator-sdk completion bash > /etc/bash_completion.d/operator-sdk
source <(operator-sdk completion bash)

# 若为x86的linux则可以直接
curl --proxy http://192.168.116.189:1087 -LO https://github.com/operator-framework/operator-sdk/releases/download/v1.11.0/operator-sdk_linux_amd64
chmod +x operator-sdk_linux_amd64 && mv operator-sdk_linux_amd64 /usr/local/bin/operator-sdk
```

#### 项目初始化

```shell
mkdir -p $GOPATH/src/github.com/HenryGits/hadoop-operator && cd $GOPATH/src/github.com/HenryGits/hadoop-operator
operator-sdk init --domain dameng.com --repo github.com/HenryGits/hadoop-operator
operator-sdk edit --multigroup=true
operator-sdk create api --group hadoop --version v1 --kind Hadoop --resource --controller
```

#### code-generator 的使用
[kubernetes/code-generator: Generators for kube-like API types](https://github.com/kubernetes/code-generator)

[Kubernetes Deep Dive: Code Generation for CustomResources](https://cloud.redhat.com/blog/kubernetes-deep-dive-code-generation-customresources)

[万物皆可operator之三，code generator的补充](https://blog.csdn.net/sixinchao_1/article/details/109997736)

```shell
echo -e "
update-codegen:
  @./hack/update-codegen.sh " >> Makefile
```
1. 执行 `go get k8s.io/code-generator@v0.22.1` 下载 code-generator 到 `$GOPATH/pkg/mod/k8s.io/`
2. 将项目拷贝到 `$GOPATH/src/github.com/HenryGits/hadoop-operator` 目录下（必须，不能将项目直接放置到 src 目录下）
3. 执行 `make update-codegen` 即可生成 clientset、informers、listers 代码


5、安装 Kubebuilder

```shell
curl --proxy http://192.168.116.189:1087 -L -o kubebuilder https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)
chmod +x kubebuilder && mv kubebuilder /usr/local/bin/
echo 'source <(kubebuilder completion bash)' >>~/.bashrc
kubebuilder completion bash > /etc/bash_completion.d/kubebuilder
source <(kubebuilder completion bash)
```

#### 使用说明

1. make generate && make manifests && make
2. make docker-build && make docker-push
3. make deploy
4. kubectl apply -f config/samples/hadoop_v1_hadoop.yaml
5. kubectl delete -f config/samples/hadoop_v1_hadoop.yaml
6. make undeploy

```
如果您尝试将 CustomResources 与基于 Kubernetes 1.8 的 client-go 一起使用——有些人可能已经很高兴了，因为他们不小心提供了 master 分支的 k8s.op/apimachinery——你遇到了 CustomResource 类型所做的编译器错误未实现

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
```

#### 目录布局：

https://access.redhat.com/documentation/zh-cn/openshift_container_platform/4.7/html/operators/go-based-operators
```
docker rm `docker ps -a -q`
docker rmi hadoop:v3.3.1
# 删除为<none>的镜像
docker rmi `docker images |awk  '{print $1,$3}'|grep none|awk  '{print $2}'`


docker build -f Dockerfile-Hadoop -t hadoop:v3.3.1   .
docker build -t controller:latest --network host --build-arg HTTP_PROXY=http://192.168.101.88:3128 --build-arg HTTPS_PROXY=http://192.168.101.88:3128 .
```

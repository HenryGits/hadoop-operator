FROM centos:7
LABEL REPOSITORY=zhc/hadoop TAG=3.3.1

#RUN echo $' \n\
#[base]  \n\
#name=CentOS7-Base \n\
#baseurl=http://mirrors.163.com/centos/7/os/x86_64/ \n\
#gpgcheck=0 \n\
#gpgkey=http://mirrors.163.com/centos/RPM-GPG-KEY-CentOS-7' > /etc/yum.repos.d/CentOS-Base.repo
#RUN cat /etc/yum.repos.d/CentOS-Base.repo
#RUN yum clean all && yum makecache fast
#
#RUN yum install -y --nogpgcheck net-tools

# jdk依赖包
ENV JAVA_HOME=/usr/local/jdk
ENV CLASSPATH=.:${JAVA_HOME}/jre/lib/rt.jar:${JAVA_HOME}/lib/dt.jar:${JAVA_HOME}/lib/tools.jar

# Hadoop 3.x独有配置
ENV HDFS_NAMENODE_USER=root
ENV HDFS_DATANODE_USER=root
ENV HDFS_SECONDARYNAMENODE_USER=root
ENV YARN_RESOURCEMANAGER_USER=root
ENV YARN_NODEMANAGER_USER=root
ENV HADOOP_HOME=/usr/local/hadoop
ENV HADOOP_CONF_DIR=$HADOOP_HOME/etc/hadoop
ENV HADOOP_CLASSPATH=$HADOOP_HOME/share/hadoop/tools/lib/*:$HADOOP_CLASSPATH
ENV PATH=$PATH:${JAVA_HOME}/bin:$HADOOP_HOME/bin:$HADOOP_HOME/sbin


COPY  jdk /usr/local/jdk
COPY  hadoop /usr/local/hadoop
COPY  Hadoop /Hadoop

RUN chmod +x /Hadoop

EXPOSE 9870 9871 8088 19888

ENTRYPOINT ["/Hadoop"]

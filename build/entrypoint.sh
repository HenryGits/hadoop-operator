#!/bin/bash

$HADOOP_HOME/bin/hadoop version


echo $HDFS_NAME_NODE_DIR

namedir=`echo $HDFS_CONF_dfs_namenode_name_dir | perl -pe 's#file://##'`


if [ -z "$CLUSTER_NAME" ]; then
  echo "Cluster name not specified"
  exit 2
fi

echo "remove lost+found from $namedir"
rm -r $namedir/lost+found


if [ ! -n $HDFS_NAME_NODE_DIR ]; then
  # 判断变量是否为空
  if [ ! -d $namedir ]; then
    echo "Namenode name directory not found: $namedir"
    exit 2
  fi

  if [ "`ls -A $HDFS_NAME_NODE_DIR`" == "" ]; then
    echo "Formatting namenode name directory: $HDFS_NAME_NODE_DIR"
    $HADOOP_HOME/bin/hdfs --config $HADOOP_CONF_DIR namenode -format $CLUSTER_NAME
  fi
fi



$HADOOP_HOME/bin/hdfs --config $HADOOP_CONF_DIR namenode

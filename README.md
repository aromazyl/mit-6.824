# mit-6.824

## Lab1: MapReduce

### Goal
实现一个Distributed MapReduce框架，该框架实现流程为:

1. 用户提供一些输入文件, map和reduce函数，以及reduce任务个数
2. Master节点启动rpc服务，等待worker节点register。Master通过调用schedule函数来安排map task给
注册过的worker节点。
3. Master把每个输入文件看做一个map task(这里和Google的paper有些区别,paper里是要把文件分成
多个map task). Master通过worker节点的DoTask RPC来调用doMap函数。doMap函数读取输入文件
的内容，并将文件名和文件内容传给用户定义的map函数，执行完后将结果按key hash到不同的中间文件。
4. 所有map task执行完后，Master继续通过调用schedule函数来安排reduce task给注册过的worker节点。
同样的，Master通过worker姐弟那的DoTask RPC来调用doReduce函数。doReduce函数读取合并对应的中间结果
，并将中间结果的k-v传递给用户定义的reduce函数，执行完后将结果输出到最终的文件里。
5. Master将map产生的n个最终文件合并成一个文件（paper里提到其实这步没有必要）
6. Master发送shutdown rpc给各个worker，然后shutdown自己

## Implemention
lab1只要我们写doMap, doReduce, schedule三个函数，比较简单。但是整个框架的代码还是很有趣的，可以都
读一读。

具体的实现直接看代码和注释。

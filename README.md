# mit-6.824

## Lab1: MapReduce
### Preparation
[MapReduce论文](https://pdos.csail.mit.edu/6.824/papers/mapreduce.pdf)

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

### Implemention
lab1只要我们写doMap, doReduce, schedule三个函数，比较简单。但是整个框架的代码还是很有趣的，可以都
读一读。

具体的实现直接看代码和注释。


## Lab2: Raft
### Preparation
[GFS论文](https://pdos.csail.mit.edu/6.824/papers/gfs.pdf)阅读整理:

* GFS设计场景
    1. System是由大量的普通计算机组成的，其中某些计算机出故障是常态。
    2. System中存储了许多大文件（multi-MB or multi-GB），小文件也有，但是不需要优化它们
    3. System通常有两种Read: large streaming reads和small random reads
    4. System通常有一种Write: large sequential write that append data to files.(small write at arbitrary当然也支持，但不优化它)
    5. System需要有效地支持多个client同时Write同一个file.
    6. 高带宽比低延迟重要
* GFS的接口
    1. create, delete, open, close, read, and write files
    2. snapshot: creates a copy of a file or a directory tree at low cost.
    3. record append: 允许多个client同时append一个文件，保证每个append都是原子操作.
* GFS Architecture
    在GFS中，一个文件被分成多个固定大小的chunks,每个chunk有唯一的64bit的chunk handle来识别。
    * Master: 在chunk创建时，给每个chunk分配chunk handle; 保存文件系统的所有metadata，包括namespace, access control information, mapping from files to chunks, current location of chunks.
    * ChunkServer: 将chunk以linux文件的形式保存在本地，通过指定chunk handle和byte range来读写文件。通常一个chunk保存在多个chunkserver。
    * client: 实现GFS的API，与master和chunkserver通信来读写数据。
* Chunk Size
    Chunk Size是一个重要的设计参数, GFS里选的是64MB。
    
    为什么要选这么大的Chunk Size 呢？
    
    好处：减少与master的通信次数；降低网络负载（不需要频繁的切换chunk，重新建立TCP连接）；减小metadata的大小

    坏处: internal fragmentation导致空间浪费；频繁使用同一个chunk导致hot spot
* Metadata
    Master主要保存三种Metadata: the file and chunk namespaces, the mapping from files to chunks, the locations of each chunk’s replicas. 
    
    Metadata的设计：
    1. 这些数据直接保存在内存中。前两种metadata也会定期持久化到磁盘。
    2. chunk location是在Master启动时向chunkserver询问的，并且通过HeartBeat定期询问。
    3. Operation log会记录metadata的change，这些log保存在master和remote（备份）。Master会定期对内存中的metadata作checkpoint，并舍弃checkpoint之前的operation log。当Master重启时，它恢复checkpoint的metadata，并执行operation log，保证metadata是最新的。

_TODO: 看不下去了 有时间继续_

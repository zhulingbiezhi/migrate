* 背景介绍：
    - 假设A表为旧表，B为新表
    - 设置redis的key：migrate_user = 1100
        - 高2位为 A表的读写标志
        - 低2位为 B表的读写标志

    - `注意：代码逻辑中默认未取到redis的migrate_user 结果时，set migrate_user = 1100`

* 定时任务：每隔3秒迁移A表1000条数据到B表，每次都获取redis的migrate_user_end的值，每次取ID的上限为这个值
* 迁移步骤：
    - 1、上线新版本代码，等待多节点完全上线，此时读A、写A
    - 2、打开写B表的开关，此时读A、写A、写B，即set migrate_user = 1101
        - ``
        写入A表时，同时获取insert ID，将数据插入到B表，
        当第一条数据记录到B表时，利用redis的setnx特性，设置此时的insert ID，set migrate_user_end = ID
        ``
    - 3、等待定时任务完成，即A，B表的数据一致
    - 4、开启读B，关闭读A，此时读B、写A、写B，即set migrate_user = 0111，观察异常
    - 5、关闭写A，此时读B、写B，即set migrate_user = 0011

* 回滚机制：
    - 步骤1,2,3  --> 直接代码回滚
    - 步骤4 --> set migrate_user = 1101
    - 步骤5 --> set migrate_user = 1101
    
* 注意：
    - 在定时任务迁移进度达到90%以上再进行上线步骤
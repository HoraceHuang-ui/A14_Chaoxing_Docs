# 接口
## 学生端
### 总评分级
```json
{
  "userid": "账号",
  "name": "xxx",
  "clazz": "班级编号",
  "grade": 0 // 0/1/2/3 代表四个总评分级
}
```
### 词云
```json
{
  "userid": "账号",
  "name": "xxx",
  "clazz": "班级编号",
  "tags": [
    {
      "tagId": 0,
      "tagName": "标签名字",
      "tagDesc": "标签描述",
      "importance": 0.34 // (0,1] 重要性越大则在词云中这个词显示越大
    },
    {
      // ...
    }
  ]
}
```
### 数据块
```json
{
  "userid": "账号",
  "name": "xxx",
  "clazz": "班级编号",
  "blocks": [
    {
      "blockId": 0, // 0/1/2，共三个数据块，数据块含义前端维护应该可以
      "blockVal": "0" // 用字符串是因为有一个数据块要显示百分比
    },
    {
      // ...
    }
  ]
}
```
### 饼图
```json
{
  "userid": "账号",
  "name": "xxx",
  "clazz": "班级编号",
  "graphs": [
    {
      "title": "上课出勤率",
      "dataVals": [0,2,20,1],
      "dataDesc": ["缺勤","迟到","正常","请假"]
    },
    {
      "title": "作业情况",
      "dataVals": [2,8,6,0],
      "dataDesc": ["满分","优秀","合格","不合格"]
    }
  ]
}
```
### 雷达图
```json
{
  "userid": "账号",
  "name": "xxx",
  "clazz": "班级编号",
  "dataVals": [4, 5, 5, 4, 4], // n个0-6的整数
  "dataDesc": ["上课态度","学习主动性","作业情况","学习时长","考试情况"]
}
```

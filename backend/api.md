# API文档

## Auth

- POST /auth/login，用户登录

  - 请求体：

    ```json
    {
        "username":"",
        "password":""
    }
    ```

  - 响应：

    ```json
    {
    "message": "登录成功/失败",
    "uid": ,
    "usertype": 
    }
    ```

- POST /auth/register，用户注册

  - 请求体：

    ```json
    {
        "Username":"",
        "Password":"", 
        "Name":"",
        "Gender":"",
        "Phone":"",
        "Address":"" //允许为null
    }
    ```

  - 响应：
  
    ```json
    {
        "message": "注册成功/失败"
    }
    ```

## User

- GET /user/info，获取用户信息

  - 请求体

  ```json
  无
  ```

  - 响应体
  
  ```json
    {
        "data": {
            "username": "",
            "name": "",
            "gender": "F/M",
            "phone": "",
            "address": "",
            "avatar": "/image/avatar/default/rice.jpg" // 这是默认图片
        },
        "message": "User information retrieved successfully" // 或者是报错
    }
  ```

- PATCH /user/update，更新用户信息

  - 请求体

  ```json
    {
        "username": "",
        "name": "",
        "gender":"F/M",
        "phone":"",
        "address":""
    }
  ```

  - 响应体
  
  ```json
  {
      "message": "User information updated successfully"
  }
  ```

- PATCH /user/avatar，上传头像

  - 请求体

  ```
  multipart/form-data //文件类型
  avatar: 图片文件 // key:value
  ```

  - 响应体
  
  ```json
  {
      "message": "User avatar updated successfully",
      "avatar_url": "/image/avatar/user/xxx.jpg"
  }
  ```

- DELETE /user/delete,软删除用户

  - 请求体

  ```json
    无
  ```

  - 响应体
  
  ```json
  {
      "message": "User deleted successfully"
  }
  ```

## Food

- POST /food/add，添加食物

  - 请求体

  ```
  multipart/form-data //文件类型
  name:""
  description:""
  image:图片文件（file）
  ```

  - 响应体
  
  ```json
  {
      "message": "Food added successfully"
  }
  ```

- DELETE /food/delete，删除已有食材（根据name）

  - 请求体

  ```json
    {
        "name":"" //食材name
    }
  ```

  - 响应体
  
  ```json
  {
    "message": "Food deleted successfully"
  }
  ```

- GET /food/all, 获取所有食材

  - 请求体

  ```json
    无
  ```

  - 响应体
  
  ```json
    {
        "foods": [
            {
                "description": "",
                "image": "", //图片的url
                "name": ""
            }]
    }
  ```

- GET /food/get/:name，根据名字获取食材信息

  - 请求体

  ```josn
    {
        "name":"水"
    }
  ```

  - 响应体
  
  ```json
  {
      "message": "Food added successfully"
  }
  ```

- PATCH /food/update, 更新食材信息

  - 请求体

  ```
  multipart/form-data //文件类型
  name:""
  description:""
  image:图片文件（file）
  ```

  - 响应体
  
  ```json
  {
      "message": "Food added successfully"
  }
  ```

## Recipe

- /recipe/list，获取菜谱列表

  - 请求体

  ```json
  无
  ```

  - 响应体
  
  ```json
  {
      "data": [
          {
              "id": 1,
              "name": "",
              "ingredients": "",
              "steps": "",
              "difficulty": "",
              "time": "",
              "author": "",
              "images": ["/image/recipe/xxx/1.jpg"]
          }
      ],
      "message": "获取菜谱列表成功/失败"
  }
  ```

- /recipe/add，添加菜谱

  - 请求体

  ```json
  {
      "name": "",
      "ingredients": "",
      "steps": "",
      "difficulty": "",
      "time": ""
  }
  ```

  - 响应体
  
  ```json
  {
      "message": "添加菜谱成功/失败",
      "recipe_id": 1
  }
  ```

- /recipe/upload，上传菜谱图片

  - 请求体

  ```
  multipart/form-data
  recipe_images: 图片文件(可多个)
  recipe_id: 菜谱ID
  ```

  - 响应体
  
  ```json
  {
      "message": "菜谱图片上传成功/失败",
      "image_urls": ["/image/recipe/xxx/1.jpg"]
  }
  ```

## Label

- /label/list，获取标签列表

  - 请求体

  ```json
  无
  ```

  - 响应体
  
  ```json
  {
      "data": [
          {
              "id": 1,
              "name": "",
              "description": ""
          }
      ],
      "message": "获取标签列表成功/失败"
  }
  ```

- /label/add，添加标签

  - 请求体

  ```json
  {
      "name": "",
      "description": ""
  }
  ```

  - 响应体
  
  ```json
  {
      "message": "添加标签成功/失败",
      "label_id": 1
  }
  ```

## Comment

- /comment/list，获取评论列表

  - 请求体

  ```json
  {
      "recipe_id": 1,
      "page": 1,
      "page_size": 10
  }
  ```

  - 响应体
  
  ```json
  {
      "data": [
          {
              "id": 1,
              "user_id": 1,
              "username": "",
              "content": "",
              "created_at": "",
              "likes": 0
          }
      ],
      "total": 100,
      "message": "获取评论列表成功/失败"
  }
  ```

- /comment/add，添加评论

  - 请求体

  ```json
  {
      "recipe_id": 1,
      "content": ""
  }
  ```

  - 响应体
  
  ```json
  {
      "message": "添加评论成功/失败",
      "comment_id": 1
  }
  ```

- /comment/like，点赞评论

  - 请求体

  ```json
  {
      "comment_id": 1
  }
  ```

  - 响应体
  
  ```json
  {
      "message": "点赞成功/取消点赞成功/失败"
  }
  ```

## Seniority

- /seniority/list，获取资历列表

  - 请求体

  ```json
  无
  ```

  - 响应体
  
  ```json
  {
      "data": [
          {
              "id": 1,
              "user_id": 1,
              "username": "",
              "level": "",
              "experience": 0,
              "badges": []
          }
      ],
      "message": "获取资历列表成功/失败"
  }
  ```

- /seniority/update，更新用户资历

  - 请求体

  ```json
  {
      "experience": 10
  }
  ```

  - 响应体
  
  ```json
  {
      "message": "更新资历成功/失败",
      "current_level": "",
      "current_experience": 100
  }
  ```

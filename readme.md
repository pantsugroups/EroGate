EroGate ---- Ero的網關入口
-------
使用了Echo框架

詳細文檔地址請訪問：https://echo.labstack.com/
github地址：https://github.com/labstack/echo

請在`gateway/conf.d`下面動態添加yaml文件。這個作爲路由使用，詳細結構如下：
```yaml
route: /balabala/*
backend: https://localhost:8080/balabala
```

## HOT TO USE?
編譯：
```bash
./install.sh || ./install.bat
```

<del>

gateway/ ----- <del>注冊的中間件(middleware)</del>Route處理函數
-------
<del>中間件用途：
    獲取訪問的http報文header中的x-headers-session值並給與驗證
    如果通過則把數據按協議打包並轉發給注冊后的不同後端地址</del>
    
中間件個錘錘，看了下代碼基本等於重新實現了一遍Route。
直接都用這個處理函數就行了



config/ ----- 用於讀取配置文件
-------
  |
  -- config.exe 添加配置文件的cli
  |
  -- config.yml 配置文件地址
</del>

整個錘錘那麽多花里胡哨的。

我一個文件就夠了daze。

DIO我不做人啦~~~~~~
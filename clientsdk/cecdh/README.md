# mac 下的交叉编译
```
./config --prefix=/opt/openssl shared no-asm 
make
cp -rf ./libcrypto.a ../bin/libcrypto_darwin.a
make clean
./Configure --prefix=/opt/openssl --cross-compile-prefix=x86_64-w64-mingw32- shared no-asm mingw64
make
cp -rf ./libcrypto.dll.a ../bin/libcrypto_windows.dll.a
make clean
./Configure --prefix=/opt/openssl --cross-compile-prefix=x86_64-linux-musl- shared no-asm linux-x86_64
make
cp -rf ./libcrypto.a ../bin/libcrypto_linux.a
make clean
```

windows:
brew install mingw-w64
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build main.go

linux:
brew install FiloSottile/musl-cross/musl-cross
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CGO_LDFLAGS="-static" go build -a main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
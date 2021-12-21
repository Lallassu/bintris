# bintris
Binary Tetris

## OpenAL
1. Update toolchain.cmake to use `-O3 -s` to build smaller version of OpenAL.
2. Update toolchain.cmake in the NDK with `-m64` to build 64bit library and some optimization to keep library size down.
```
list(APPEND ANDROID_COMPILER_FLAGS "-m64")
list(APPEND ANDROID_COMPILER_FLAGS_RELEASE -O -W1 -s -w -g0)
```

### Changes to gomobile to compile 64bit lib
cmd := exec.Command(cmake, "-S",
        initOpenAL,
        "-DANDROID_PLATFORM=23",
        "-B", buildDir,
        //"-DCMAKE_TOOLCHAIN_FILE="+initOpenAL+"/XCompile-Android.txt",
        "-DCMAKE_TOOLCHAIN_FILE="+ndkRoot+"/build/cmake/android.toolchain.cmake",
        "-DNDEBUG",
        "-DANDROID_HOST_TAG="+t.ClangPrefix())

## Optimize Build Size
1. Use `ldflags="-w" to remove debug information from the build

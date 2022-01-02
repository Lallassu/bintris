# bintris
Binary Tetris

## OpenAL
1. Update toolchain.cmake to use `-O3 -s` to build smaller version of OpenAL.
2. Update toolchain.cmake in the NDK with `-m64` to build 64bit library and some optimization to keep library size down.
```
list(APPEND ANDROID_COMPILER_FLAGS "-m64")
list(APPEND ANDROID_COMPILER_FLAGS_RELEASE -O -W1 -s -w -g0)
```

### Changes to gomobile cmd
Line 182 in `mobile/cmd/init.go`:
```go
		cmd := exec.Command(cmake, "-S",
			initOpenAL,
			"-DANDROID_PLATFORM=23",
			"-B", buildDir,
			"-DCMAKE_TOOLCHAIN_FILE="+ndkRoot+"/build/cmake/android.toolchain.cmake",
			"-DANDROID_HOST_TAG="+t.ClangPrefix())
```

## Optimize Build Size
1. Use `ldflags="-w" to remove debug information from the build

## audio/al updates
Using AAB packaging for PlayStore requires libraries loaded with `dlopen` to not use path. Hence,
the audio/al packages requires to just use `dlopen("libopenal.so")`. The base.apk in aab doesn't include
the libraries rather they are included in the `split_config.<arch>.apk`.

AL/openal.h etc requires to be in audio/al package dir.

exp/audio/al/al_android.go:17 
```c
 *handle = dlopen("libopenal.so", RTLD_LAZY);
```

Debug using:
`adb shell pm list packages -f |grep bintris`
Then download:
`adb pull <path to bintris package base path>`

## android.toolchain.cmake
Optimize options and also armv8 (line 578)


## Build
1. Build bintris `make android`. 
2. Unzip bintris.apk
3. Move lib to `jniLibs` in Android Studio project.
4. Build 

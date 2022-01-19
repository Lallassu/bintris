# 🎮 Bintris
<img src=https://github.com/lallassu/bintris/raw/main/raw_assets/screenshots/front.jpg width=20%/> <img src=https://github.com/lallassu/bintris/raw/main/raw_assets/screenshots/play.jpg width=20%/> <img src=https://github.com/lallassu/bintris/raw/main/raw_assets/screenshots/bitrot.jpg width=20%/> <img src=https://github.com/lallassu/bintris/raw/main/raw_assets/screenshots/scoreboard.jpg width=20%/> 

<h3>Bintris is a mobile game developed in Go!</h3>

Support me by buying this game on Google PlayStore 💸: https://play.google.com/store/apps/details?id=go.lang.bintris

Or you can of course build it from source yourself (see instructions below 🪛).

Enjoy!

## About The Game
Bintris is a small game inspired by Tetris. The goal is to flip the bits so that the
bits represent the decimal number in the right column. When the bits represent the decimal number
the line is cleared and points are gathered. Number of bits representing the decimal number is also how
much points scored for the particular line.

## Demo On Youtube
[![Bintris](https://img.youtube.com/vi/nK01rQ2STPs/0.jpg)](https://www.youtube.com/watch?v=nK01rQ2STPs)


## About The Implementation
The game is developed in Go and is implemented using OpenGL (graphics) and OpenAL (sound). Gomobile is used
to generate shared libraries that are used for the Android build. The game works just as good on Linux as on Android.

It all started as an experiment with Gomobile and ended up as a fully working game, after a lot of frustration and gotchas! ;)
With that said, the source is a bit of a mess and I have some stuff on my todo-list such as implementing full
parsing of wav header, a simple FSM for the game etc.


## Building From Source

To just run the game on Linux, just issue the commands:
```bash
# First remove line for modified version of gomobile
sed -i '/nergal/d' go.mod
# Then run the game
go run .
```

## Build For Android
First make sure you have built OpenAL (`make openal`). It's a bit complicated to make the OpenAL
build work succesfully with current state of `gomobile`. I've added some notes below about this.

```
make android
```

## Build For Android Studio
The repository includes a very small Android Studio project that is used to build AAB package format for Google PlayStore. This
project handles AAB packaging and uses the shared objects (.so) files from the build.

To build/run via Android Studio (make sure to have OpenAL libraries first `make openal`, see requirements below):

1. `make studio`
2. Open the project in Android Studio
3. Attach mobile or virtual device and run.

## Details

### OpenAL
Update toolchain.cmake to use `-O3 -s` to build smaller version of OpenAL (otherwise you will end up with ~30MB debug none stripped version)
Also update to armv8 (line 578) in the toolchain from v7.

Using AAB packaging for PlayStore requires libraries loaded with `dlopen` to not use the path. Hence,
the audio/al packages requires to just use `dlopen("libopenal.so")`. The base.apk in aab doesn't include
the libraries rather they are included in the `split_config.<arch>.apk`.

AL/openal.h etc requires to be in audio/al package dir for rebuilding:

exp/audio/al/al_android.go:17 
```c
 // All other code for reading ENV can be removed.
 *handle = dlopen("libopenal.so", RTLD_LAZY);
```

Build OpenAL:

```
make openal
```

### Go
Use `ldflags="-w" to remove debug information from the build (see Makefile)

#### Changes Required to Gomobile Command
Line 182 in `mobile/cmd/init.go`:
```go
cmd := exec.Command(cmake, "-S",
	initOpenAL,
	"-DANDROID_PLATFORM=23",
	"-B", buildDir,
	"-DCMAKE_TOOLCHAIN_FILE="+ndkRoot+"/build/cmake/android.toolchain.cmake",
	"-DANDROID_HOST_TAG="+t.ClangPrefix())
```

## ADB Debug
View logs from connected phone (developer mode):
```
adb logs
```

Debug using:
```
adb shell pm list packages -f |grep bintris
```

Then download:
```
adb pull <path to bintris package base path>
```

Too see what is included in the divided apk's.

## License
GNU General Public License v3.0 (see COPYING)


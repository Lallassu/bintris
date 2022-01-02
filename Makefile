NDK=~/Android/Sdk/ndk/23.0.7599858
AL_PATH=~/Downloads/openal/openal-soft/
APK_TMP=/tmp/_bintris_apk_/
ANDROID_HOME=~/Android/Sdk/
CMAKE_PATH=~/Downloads/cmake-3.22.1-linux-x86_64/bin
GOMOBILE_PATH=~/go/src/golang.org/x/mobile/cmd/gomobile/gomobile

# Make studio builds shared objects and places them in the AS project dir.
studio: openal android
	unzip bintris.apk -d $(APK_TMP)
	cp -r $(APK_TMP)/lib/* AndroidProject/app/src/main/jniLibs/

openal:
	PATH="$(CMAKE_PATH):$(PATH)" ANDROID_PLATFORM=android-23 ANDROID_HOME=$(ANDROID_HOME) ANDROID_NDK_HOME=$(NDK) $(GOMOBILE_PATH) init -openal $(AL_PATH)

init:
	ANDROID_NDK_HOME=$(NDK) gomobile init -openal $(AL_PATH) -ldflags="-w" -gcflags="-w -I $(AL_PATH)/include/" 

android:
	ANDROID_NDK_HOME=$(NDK) $(GOMOBILE_PATH) build -target android -androidapi=23 -o bintris.apk -ldflags="-w" -gcflags="-w -I $(AL_PATH)/include/" 

bind:
	ANDROID_HOME=$(NDK) ANDROID_NDK_HOME=$(NDK) gomobile bind -target android -androidapi=23 -o bintris.apk -ldflags="-w" -gcflags="-w -I $(AL_PATH)/include/" 

android_install:
	ANDROID_NDK_HOME=$(NDK) $(GOMOBILE_PATH) install -target android -androidapi=23 -o bintris.apk -ldflags="-w" -gcflags="-I $(AL_PATH)/include/" -x

linux:
	go build .

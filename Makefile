#NDK=~/Android/Sdk/ndk-bundle
NDK=~/Android/Sdk/ndk/23.0.7599858

openal:
	#ANDROID_HOME=~/Android/Sdk/ ANDROID_NDK_HOME=$(NDK) /home/nergal/go/src/golang.org/x/mobile/cmd/gomobile/gomobile init -openal /home/nergal/Downloads/openal/openal-soft/ -v -x
	#PATH="/home/nergal/Downloads/cmake-3.22.1-linux-x86_64/bin:$(PATH)" ANDROID_PLATFORM=android-23 ANDROID_HOME=~/Android/Sdk/ ANDROID_NDK_HOME=$(NDK) gomobile init -openal /home/nergal/Downloads/openal/openal-soft/ -v -x
	#PATH="/home/nergal/Downloads/cmake-3.22.1-linux-x86_64/bin:$(PATH)" ANDROID_PLATFORM=android-23 ANDROID_HOME=~/Android/Sdk/ ANDROID_NDK_HOME=$(NDK) /home/nergal/go/src/golang.org/x/mobile/cmd/gomobile/gomobile init -openal /home/nergal/Downloads/openal/openal-soft/ -v -x
	PATH="/home/nergal/Downloads/cmake-3.22.1-linux-x86_64/bin:$(PATH)" ANDROID_PLATFORM=android-23 ANDROID_HOME=~/Android/Sdk/ ANDROID_NDK_HOME=$(NDK) /home/nergal/go/src/golang.org/x/mobile/cmd/gomobile/gomobile init -openal /home/nergal/Downloads/openal/openal-soft/ -v -x
android:
	#ANDROID_NDK_HOME=$(NDK) gomobile build -target android -androidapi=23 -o bintris.apk
	ANDROID_NDK_HOME=$(NDK) /home/nergal/go/src/golang.org/x/mobile/cmd/gomobile/gomobile build -target android -androidapi=23 -o bintris.apk -gcflags="-I /home/nergal/Downloads/openal/openal-soft/include/" -x
android_install:
	#ANDROID_NDK_HOME=$(NDK) gomobile install -target android -androidapi=23 -o bintris.apk
	ANDROID_NDK_HOME=$(NDK) /home/nergal/go/src/golang.org/x/mobile/cmd/gomobile/gomobile install -target android -androidapi=23 -o bintris.apk -gcflags="-I /home/nergal/Downloads/openal/openal-soft/include/" -x
linux:
	go build .

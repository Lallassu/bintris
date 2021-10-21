NDK=~/Android/Sdk/ndk-bundle
android:
	ANDROID_NDK_HOME=$(NDK) gomobile build -target android -androidapi=23 -o bintris.apk
android_install:
	ANDROID_NDK_HOME=$(NDK) gomobile install -target android -androidapi=23 -o bintris.apk
linux:
	go build .

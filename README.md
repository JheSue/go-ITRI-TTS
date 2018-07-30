go-ITRI-TTS
===

Sign up Account:
		http://tts.itri.org.tw/member/registeration.php?alry=1

env:

	gvm use go1.9

build:

	go build .

run:

	./soap-client <ACCOUNT> <PASSWPRD> <TEXT>

example:

        ./soap-client <ACCOUNT> <PASSWORD> "Hello I am Jarvis"

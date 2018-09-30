build:
	docker build -t stl10n/martian .

push:
	docker push stl10n/martian

deploy:
	go build .
	scp martian st.gortc.io:/usr/bin/martian
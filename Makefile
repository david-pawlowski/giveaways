.PHONY: deploy
deploy:
	docker build --platform linux/amd64 -t karamba116/giveaways:azure .
	docker push karamba116/giveaways:azure

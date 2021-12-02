test:
	shellcheck testdata/*.sh
	./bin/shreq testdata/*.sh

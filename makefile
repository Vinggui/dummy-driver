# This will install the necessary tools for this driver and allows to run
.PHONY: install run
install:
	@go build main.go

run:
	@./main $(TOKEN)

BOLD := $(shell tput -T linux bold)
GREEN := $(shell tput -T linux setaf 2)
PURPLE := $(shell tput -T linux setaf 5)
RESET := $(shell tput -T linux sgr0)

SUCCESS := $(BOLD)$(GREEN)
TITLE := $(BOLD)$(PURPLE)

.PHONY: run start stop
run:
	@go run .

start:
	@printf '$(TITLE)Building now...$(RESET)\n'
	@go mod tidy & go build -o zhtcloud .
	@nohup ./zhtcloud > run.log 2>&1 &
	@printf '$(SUCCESS)Build and start successfully!\n$(RESET)'

stop:
	@pkill zhtcloud
	@rm -f run.log
	@printf '$(SUCCESS)Clean successfully!\n$(RESET)'
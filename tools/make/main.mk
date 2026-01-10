TAIL__MAKE_PROJECT_PATH := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
include $(TAIL__MAKE_PROJECT_PATH)/internal/docker.mk
include $(TAIL__MAKE_PROJECT_PATH)/lang/ru.mk

TAIL__ROOT_PATH := $(CURDIR)
TAIL__APP_PATH = "$(TAIL__ROOT_PATH)/.temp"

include $(TAIL__MAKE_PROJECT_PATH)/tests/integration.mk
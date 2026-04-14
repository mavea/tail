TAIL__LANG_PREFIX := $(shell printf '%s' "$(LANG)" | tr '[:upper:]' '[:lower:]' | cut -c1-2)
ifeq ($(TAIL__LANG_PREFIX),ru)
include $(TAIL__MAKE_PROJECT_PATH)/lang/ru.mk
else ifeq ($(TAIL__LANG_PREFIX),en)
include $(TAIL__MAKE_PROJECT_PATH)/lang/en.mk
else
include $(TAIL__MAKE_PROJECT_PATH)/lang/en.mk
endif
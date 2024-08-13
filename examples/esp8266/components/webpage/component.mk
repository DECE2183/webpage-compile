COMPONENT_EXTRA_CLEAN := $(COMPONENT_PATH)/pages/bin/*
COMPONENT_EMBED_FILES := pages/bin/pages.bin

webpage.o: $(COMPONENT_PATH)/pages/bin/pages.bin

$(COMPONENT_PATH)/pages/bin/pages.bin:
	webpage-compile $(COMPONENT_PATH)/pages/pages.yaml

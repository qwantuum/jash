use zed_extension_api as zed;

struct JashExtension;

impl zed::Extension for JashExtension {
    fn new() -> Self {
        Self
    }
}

zed::register_extension!(JashExtension);

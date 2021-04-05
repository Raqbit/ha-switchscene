# ha-switchscene

Small utility which I use to activate Home Assistant scenes via keyboard shortcuts.

The utility will also send notifications when it changs the scene or when errors occur, although this currently only
works on Linux-based systems with the `notify-send` application present. Install `libnotify-bin` on debian or Ubuntu.

## Usage

```
ha-switchscene -url "http://homeassistant.local:8123" -scene "scene.some_scene" -name "Some Scene"
```

## Storing long-lived token

`ha-switchscene` retrieves your long-lived Home Assistant token from your keyring. To install your Home Assistant token
into your keyring, please run `ha-switchscene` in `storeToken`
mode:

```
ha-switchscene -storeToken -url "http://homeassistant.local:8123"
```

Paste in your token: Note that your input will not be shown back to you.

Every time you now use `ha-swichscene` with the same Home Assistant URL, your stored token will be used.

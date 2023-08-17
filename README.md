
<p align="center">
  <img width="300" alt="image" src="https://user-images.githubusercontent.com/13785588/179449532-a4beb00f-0315-4386-9468-e494fc347224.png">
</p>

Lunii Admin is evolving!

The new, no installation required lunii-admin-web is even better. All the great features of lunii-admin but right where and when you need it, in your browser. An it has better pack installation feedback and even some tetris going on!

Don't hesitate to check out the [hosted version](https://lunii-admin-web.pages.dev) or the [github repo](https://github.com/olup/lunii-admin-web)

And to complete this new tool, you have a new pack creator: [lunii admin builder](https://lunii-admin-builder.pages.dev). It features a nice and simple UI to build your packs in a flash, and generate your missing images or audios on the fly.

------

[Download latest version](https://github.com/olup/lunii-admin/releases)

This is a small desktop app to manage your lunii device.

Features :
- list installed story packs
- reorder story packs
- delete story pack
- Install third party packs (STUdio zip pack, more info on https://github.com/marian-m12l/studio)
- create your own thrid party story packs from a structured directory (details below)

## Creating your own story packs
https://github.com/marian-m12l/studio offers a full editor to create your story packs, by I find that overkill for most users.

This software lets you create a STUdio compatible story pack from a structured directory.

The structure needs to be as follows:
```
pack-name/
    title.mp3
    cover.jpeg
    md.yaml

    first-story/
        title.mp3
        cover.jpeg
        story.mp3
    
    second-story/
        title.mp3
        cover.jpeg
        story.mp3
    
    sub-menu/
        title.mp3
        cover.jpeg

        sub-story/
            title.mp3
            cover.jpeg
            story.mp3
```

the `md.yaml` file should have the following content:

```
title: <Your pack title>
description: <Your pack description>
uuid: <A single uuid - generate it from https://www.uuidtools.com/v4>
```

The use the "create pack" button from the app to convert such a directory to a `.zip` STUdio story pack.

You can then install the resulting pack to your device with the "install pack" button.

---
## Build from sources

Commands to build can be seen in the [GitHub Pipeline configuration file](.github/workflows/build-version.yaml)

Some dependencies should be installed first (eg: libmp3lame-dev).

## Technologies

Build with GO and React, desktop app powered by [Wails](https://github.com/wailsapp/)


# wallpaper-engine-downloader

If you can't just add wallpapers from steamworkshop, you have to download them via steamworhshop downloaders, unzip and put in projects folder.
Tired of doing this every time? Just use wallpaper-engine-downloader

## Usage

Download ready-to-go windows .exe file

### Prepare config.toml file

Locate wallpaper engine folder.
It contains directory 'projects', which contains 'myprojects'.
Copy path to 'myprojects' folder 

In the directory containing executable, create file named config.toml.
Right click, Edit

```toml
MyprojectPath = "C:\\Escaped\\Path\\To\\myproject\\folder"
```

### Using precompiled .exe file

Open console and proceed into directory containing executable

-config config-file-name.toml (should be in same directory)

-wallpaper link-to-wallpaper in steamworkshop

```bash
./wallpaper-engine-downloader.exe -config config.toml -wallpaper "https://steamcommunity.com/sharedfiles/filedetails/?id=818603284&searchtext="
```

## Credentials

This program uses http://steamworkshop.download/ to download items from steamworkshop. Future versions could support steam API to download directly from steam, but for now, I leave it as is.

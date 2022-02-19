var vlc = Application('VLC');

function run(argv) {
    let json = {
        "playing": vlc.playing(),
        "nameOfCurrentItem": vlc.nameOfCurrentItem(),
        "fullscreenMode": vlc.fullscreenMode(),
        "durationOfCurrentItem": vlc.durationOfCurrentItem(),
        "currentTime": vlc.currentTime()
    };
    return JSON.stringify(json);
}
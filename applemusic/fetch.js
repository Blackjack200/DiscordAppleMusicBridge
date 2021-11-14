const APP = Application.currentApplication();
APP.includeStandardAdditions = true;
ObjC.import('Foundation');

const iTunesApp = Application('Music');

function run(argv) {
    let track = iTunesApp.currentTrack;
    let data = {
        "name": track.name(),
        "kind": track.kind(),
        "album": track.album(),
        "artist": track.artist(),
        "bitRate": track.bitRate(),
        "discCount": track.discCount(),
        "discNumber": track.discNumber(),
        "duration": track.duration(),
        "sampleRate": track.sampleRate(),
        "trackCount": track.trackCount(),
        "trackNumber": track.trackNumber(),
    };
    return JSON.stringify(data);
}
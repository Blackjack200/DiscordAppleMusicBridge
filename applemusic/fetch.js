var music = Application('Music');

function run(argv) {
    let track = music.currentTrack;
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
        "playerState": music.playerState(),
        "playerPosition": music.playerPosition(),
    };
    return JSON.stringify(data);
}
## discord-music-link-converter

This is a discord bot that converts links between popular music services.

So far:
  - spotify track => apple track conversion is done
  - apple album recognition and GETing works
  - apple track GETing is a little weird, because they have two link formats (a canonical one, and a common one):
    - canonical => something like `https://music.apple.com/us/song/<song-name>/<song-id>`
    - common (calling it common because that's what the Apple Music website links when you try to share) => `https://music.apple.com/us/album/<album-name>/<album-id>?i=<song-id>`

In addition to finishing the above features, support for YouTube, probably using the official YouTube go client, is also planned, but hasn't been looked into it that much yet.

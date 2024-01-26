This branch is a modified version of Main which is designed to work with a single account. The refresh token is stored in a database on the server and this allows for it to remain indefinitely logged in. This is required to demo the website as I currently do not have a full access Spotify API key

https://github.com/ejagombar/SpannerBackend/assets/77460324/0cd6d9d1-2a8b-443a-ad08-6397535586c6

# Spanner

Spanner is a Spotify Analyser website. It allows you to view your top tracks and artists and analyse playlists.
Spanner is currently undergoing the verification process with Spotify to access the Public API. Once this is approved, the website can be made live.

This Repository hosts the backend code, written in Go.

The frontend, written in Typescript, can be found here: https://github.com/ejagombar/SpannerFrontend

## Project Info
I decided to undertake this project in order to learn more about web development and APIs, as well as develop my Go knowledge and learn some Typescript and React. Although this project took longer than expected, I have learnt a log and I plan to add more features in the future.

Some things that are to be added in the future are:
 - Unit Tests
 - Lyric Sentiment Analysis
 - Metrics and Server Error Alerts
 - Playlist Comparison

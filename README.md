This project is a website that displays statistics associated with games. Currently it supports [Steam](https://store.steampowered.com/) and [Twitch](https://www.twitch.tv/) displaying Playercount and Viewer/Streamcount respectively for the games if avaliable. It allows you to search for games to retrieve statistics for using the [IGDB](https://www.igdb.com/) API search function.

#Frontend
Made with React, the project is broken down into a few main components which I will detail below.

**Header**
Simply displays the header of the website. Also displays a link to my github.

**SearchGames**
This contains a html input element which takes the name of the game the users wants to search for. It also contains a button element that passes the name the user has inputted into the 'search' function to send a request to the backend to find games based on the input. Based on what is then returned it creates a bunch of 'gameFlavouredButtons' which will be explained below.

**GameFlavouredButton**
Displays the name, thumbnail and a button to select for each game. When the button is pressed, it calls back to a function passed to it which will then pass the game id, name and thumbnail to display.

**GameDisplay**
Displays the game name thumbnail and then the information retrieved from the backend. This will be the steam player numbers and twitch viewer/stream count. If information is no avaliable in either case it will simply be ommited. It then uses a 'GameStatsDisplay' to display the individual platform statistics.

**GameStatsDisplay**
Takes a name of a platform, and an icon and then a list of fields with values to display (e.g playercount). It simply then lists them in a html unordered list.

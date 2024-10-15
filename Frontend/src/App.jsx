import React, { useState } from "react"
import Header from './header.jsx'
import GameDisplay from './gameDisplay.jsx'
import SearchGames from './searchGames.jsx'
import './App.css'
const steamGameURL = "https://store.steampowered.com/app/";
const backendRequest = "https://playeranalysisbackend2-822312071652.europe-west2.run.app/getPlayerAnalysis/";
const gameHeaderURL = "https://steamcdn-a.akamaihd.net/steam/apps/appid/header.jpg";
const gamePercentageDecimalPlaces = 3;

function App() {
  const [selectedGameID, setSelectedGameID] = useState(0);
  const [gameHeaderLink, setGameHeaderLink] = useState('');
  const [gameName, setGameName] = useState("No name set")
  const [steamPlayerData, setSteamPlayerData] = useState(null);
  const [twitchViewerData, setTwitchPlayerData] = useState(null);

  return(
    <div className="masterDiv">
      <Header/>
      <SearchGames passBackSelectedTo={showAnalysis}/>
      { selectedGameID != 0 ? <GameDisplay 
      steamData={steamPlayerData} 
      twitchData={twitchViewerData}
      gamePicURL={gameHeaderLink}
      gameName={gameName}/> : null }
    </div>
  );

  //make it so header img is passed into this function as opposed to getting it via ID
  async function showAnalysis(gameID, gameName, imageURL){
    const webAPIGameRequest = backendRequest + gameID;
    const response = await fetch(webAPIGameRequest);
    const analysisjson = await response.json();
    setPlayerStats(analysisjson);
    setGameName(gameName);
    setGameHeaderLink(imageURL);
    setSelectedGameID(gameID);
  }

  function setPlayerStats(jsonFrom) {
    console.log(jsonFrom)
    if (jsonFrom["steam"] != undefined){
      setSteamPlayerData(jsonFrom.steam);
    }
    else{
      setSteamPlayerData(null);
    }

    if (jsonFrom["twitch"] != undefined){
      setTwitchPlayerData(jsonFrom.twitch);
    }
    else{
      setTwitchPlayerData(null);
    }
    
    //const playerPercentFloat = parseFloat(jsonFrom.players_as_percentage).toFixed(2 + gamePercentageDecimalPlaces) * 100;
  }
}

export default App

import "./gameDisplay.css";
import GameStatsDisplay from "./gameStatsDisplay.jsx"

function gameDisplay({steamData, twitchData, gamePicURL, gameName}){
    return (
    <div>
        <div className="gameDisplayBack">
            <b> {gameName} </b>
            <img src={gamePicURL} className="gameHeader" alt="pic"/>
        </div>
        {twitchData != null ? <GameStatsDisplay
        platformName="Twitch" 
        platformIcon="https://banner2.cleanpng.com/20180513/xie/avccgu4ho.webp"
        fields = {[{name:"Viewers", value:twitchData.current_viewers}, 
        {name:"Streams", value:twitchData.current_streams}, 
        {name:"Average stream viewers", value:twitchData.average_viewers.toFixed(2)}]}/> : null }
        {steamData != null ? <GameStatsDisplay
        platformName="Steam" 
        platformIcon="https://upload.wikimedia.org/wikipedia/commons/8/83/Steam_icon_logo.svg"
        fields = {[{name:"Players", value:steamData.current_players}, {name:"Share of total players", value:steamData.percentage_of_total.toFixed(3) + "%"}]}/> : null }
    </div>
    );
}

export default gameDisplay
import { useState } from "react";
import GameFlavouredButton from "./gameFlavouredButton.jsx";
import "./searchGames.css"

const backendRequest = "https://playeranalysisbackend2-822312071652.europe-west2.run.app/getGames/";

function searchGames({passBackSelectedTo}){
    const [gameSearch, setGameSearch] = useState('');
    const [gamesSearchResults, setGameSearchResults] = useState([]);
    const [searched, setSearchedStatus] = useState(false);

    return (
        <div className="searchGames">
            <h1>Search for Games</h1>
            <input value={gameSearch} onChange={e => setGameSearch(e.target.value)}></input>
            <button title="Search" onClick={() => { search(gameSearch)}}>Search</button>
            {gamesSearchResults.length == 0 && searched ? <h3>No results...</h3> : null}
            {gamesSearchResults.map(searchResult => 
            (<GameFlavouredButton key={searchResult.id} gameURL={searchResult.url} gameThumb={searchResult.thumb} gameName={searchResult.name} gameID={searchResult.id} callback={passBackSelectedTo}/>))}
        </div>
    )

    async function search(term){
        const backendSearchRequest = backendRequest + term;
        const header = {
            method: 'GET',
        }
        
        let response = await fetch(backendSearchRequest, header);
        const gamesSearchjson = await response.json();
        console.log(gamesSearchjson)
        setGameSearchResults(gamesSearchjson);
        setSearchedStatus(true);
    }
}

export default searchGames
import "./gameFlavouredButton.css"

function gameFlavouredButton({gameURL, gameThumb, gameName, gameID, callback}) {
    return (
    <div className="gameFlavouredButton">
        <button className="selectGame" title={gameName} onClick={() => callback(gameID, gameName, gameThumb)}>✔</button>
        <img src={gameThumb} alt={gameName}/>
        <h6 className="gameName">{gameName}</h6>
    </div>)
}

export default gameFlavouredButton
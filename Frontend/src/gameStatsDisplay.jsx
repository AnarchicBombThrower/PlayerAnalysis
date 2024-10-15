import "./statsDisplay.css"

function gameStatsDisplay({platformIcon, platformName, fields}){
    return (
    <div className="gameStatsBack">
        <b> {platformName}: </b>
        <ul className="statsList">
            {fields.map(field => (<li>{field.value} {field.name}</li>))}
        </ul>
        <img src={platformIcon} 
        className="platformImage" alt={platformName + " icon"}
        width="40"
        height="40"/>
    </div>
    )
}

export default gameStatsDisplay